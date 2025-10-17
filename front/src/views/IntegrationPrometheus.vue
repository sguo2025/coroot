<template>
    <v-form v-model="valid" ref="form" style="max-width: 800px">
        <v-alert v-if="form.global" color="primary" outlined text>
            这个工程使用了一个全局 Prometheus 配置，无法通过 UI 修改
        </v-alert>

        <v-checkbox v-model="form.use_clickhouse" label="使用 ClickHouse 进行指标存储" class="my-2" hide-details :disabled="form.global" />
        <div class="caption mb-3">启用后，ClickHouse 将代替 Prometheus 用于指标存储。</div>

        <div class="subtitle-1">Prometheus URL</div>
        <div class="caption">根因分析先锋基于存储在 Prometheus 服务器上的遥测数据工作。</div>
        <v-text-field
            outlined
            dense
            v-model="form.url"
            :rules="[$validators.notEmpty, $validators.isUrl]"
            placeholder="https://prom.example.com:9090"
            hide-details="auto"
            class="flex-grow-1"
            single-line
            :disabled="form.global || form.use_clickhouse"
        />
        <v-checkbox
            v-model="form.tls_skip_verify"
            :disabled="!form.url.startsWith('https') || form.global || form.use_clickhouse"
            label="Skip TLS verify"
            hide-details
            class="my-2"
        />

        <v-checkbox v-model="basic_auth" label="HTTP basic auth" class="my-2" hide-details :disabled="form.global || form.use_clickhouse" />
        <div v-if="basic_auth" class="d-flex gap">
            <v-text-field
                outlined
                dense
                v-model="form.basic_auth.user"
                label="username"
                hide-details
                single-line
                :disabled="form.global || form.use_clickhouse"
            />
            <v-text-field
                v-model="form.basic_auth.password"
                label="password"
                type="password"
                outlined
                dense
                hide-details
                single-line
                :disabled="form.global || form.use_clickhouse"
            />
        </div>

        <v-checkbox v-model="custom_headers" label="Custom HTTP headers" class="my-2" hide-details :disabled="form.global || form.use_clickhouse" />
        <template v-if="custom_headers">
            <div v-for="(h, i) in form.custom_headers" :key="i" class="d-flex gap mb-2 align-center">
                <v-text-field outlined dense v-model="h.key" label="header" hide-details single-line :disabled="form.global || form.use_clickhouse" />
                <v-text-field
                    outlined
                    dense
                    v-model="h.value"
                    type="password"
                    label="value"
                    hide-details
                    single-line
                    :disabled="form.global || form.use_clickhouse"
                />
                <v-btn @click="form.custom_headers.splice(i, 1)" icon small :disabled="form.global || form.use_clickhouse">
                    <v-icon small>mdi-trash-can-outline</v-icon>
                </v-btn>
            </div>
            <v-btn color="primary" @click="form.custom_headers.push({ key: '', value: '' })" :disabled="form.global || form.use_clickhouse"
                >添加 header</v-btn
            >
        </template>

        <div class="subtitle-1 mt-3">刷新间隔</div>
        <div class="caption">
            根因分析先锋多久从 Prometheus 服务器获取一次遥测数据。值必须大于
            <a href="https://prometheus.io/docs/prometheus/latest/configuration/configuration/" target="_blank" rel="noopener noreferrer"
                ><var>scrape_interval</var></a
            > 的值。
        </div>
        <v-select v-model="form.refresh_interval" :items="refreshIntervals" outlined dense :menu-props="{ offsetY: true }" :disabled="form.global" />

        <div class="subtitle-1">Extra selector</div>
        <div class="caption">一个额外的指标选择器，将添加到每个 Prometheus 查询中 (例如 <var>{cluster="us-west-1"}</var>)</div>
        <v-text-field
            outlined
            dense
            v-model="form.extra_selector"
            :rules="[$validators.isPrometheusSelector]"
            single-line
            :disabled="form.global || form.use_clickhouse"
        />

        <div class="subtitle-1">Remote Write URL (远程写入 URL)</div>
        <div class="caption">
            如果您使用的是像 VictoriaMetrics 这样的替换 Prometheus 的工具，在集群模式下，您可能需要配置一个不同的 Remote Write
            URL. 默认情况下，根因分析先锋会在上面配置的基 URL 后面添加 <var>/api/v1/write</var>。
        </div>
        <v-text-field
            outlined
            dense
            v-model="form.remote_write_url"
            :rules="[$validators.isUrl]"
            single-line
            :disabled="form.global || form.use_clickhouse"
        />

        <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
            {{ error }}
        </v-alert>
        <v-alert v-if="message" color="green" outlined text>
            {{ message }}
        </v-alert>
        <v-btn block color="primary" @click="save" :disabled="(!valid && !form.use_clickhouse) || form.global" :loading="loading">保存</v-btn>
    </v-form>
</template>

<script>
const refreshIntervals = [
    { value: 5000, text: '5 seconds' },
    { value: 10000, text: '10 seconds' },
    { value: 15000, text: '15 seconds' },
    { value: 30000, text: '30 seconds' },
    { value: 60000, text: '60 seconds' },
];

export default {
    data() {
        return {
            form: {
                url: '',
                tls_skip_verify: false,
                basic_auth: null,
                custom_headers: [],
                refresh_interval: 0,
                extra_selector: '',
                remote_write_url: '',
                use_clickhouse: false,
            },
            basic_auth: false,
            custom_headers: true,
            valid: false,
            loading: false,
            error: '',
            message: '',
        };
    },

    mounted() {
        this.get();
    },

    watch: {
        custom_headers(v) {
            if (v && !this.form.custom_headers.length) {
                this.form.custom_headers.push({ key: '', value: '' });
            }
        },
    },

    computed: {
        refreshIntervals() {
            return refreshIntervals;
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getIntegrations('prometheus', (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.form = Object.assign({}, this.form, data);
                if (!this.form.basic_auth) {
                    this.form.basic_auth = { user: '', password: '' };
                    this.basic_auth = false;
                } else {
                    this.basic_auth = true;
                }
                if (!this.form.custom_headers) {
                    this.form.custom_headers = [];
                }
                this.custom_headers = !!this.form.custom_headers.length;
            });
        },
        save() {
            this.loading = true;
            this.error = '';
            const form = JSON.parse(JSON.stringify(this.form));
            if (!this.basic_auth) {
                form.basic_auth = null;
            }
            if (!this.custom_headers) {
                form.custom_headers = [];
            }
            this.message = '';
            this.$api.saveIntegrations('prometheus', 'save', form, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('refresh');
                this.message = 'Settings were successfully updated. The changes will take effect in a minute or two.';
                setTimeout(() => {
                    this.message = '';
                }, 3000);
                this.get();
            });
        },
    },
};
</script>

<style scoped>
.gap {
    gap: 16px;
}
</style>
