<template>
    <v-dialog v-if="check" :value="value" @input="emitValue" max-width="800">
        <v-card class="pa-5">
            <div class="d-flex align-center font-weight-medium mb-4">
                <template v-if="check.id === 'SLOAvailability' || check.id === 'SLOLatency'">
                    配置 "{{ check.title }}" 检查
                    <v-btn v-if="form && !form.default" small icon @click="confirmation = true"><v-icon small>mdi-trash-can-outline</v-icon></v-btn>
                    <v-overlay :value="confirmation" absolute opacity="0.8">
                        <div>确定要删除 "{{ check.title }}" 检查的覆盖吗？</div>
                        <div class="mt-5 d-flex">
                            <v-spacer />
                            <v-btn @click="confirmation = false" small color="info">取消</v-btn>
                            <v-btn @click="del" :loading="deleting" color="error" class="ml-3" small>删除</v-btn>
                        </div>
                    </v-overlay>
                </template>
                <template v-else> 调整 "{{ check.title }}" 检查的阈值</template>
                <v-spacer />
                <v-btn icon @click="emitValue(false)"><v-icon>mdi-close</v-icon>关闭</v-btn>
            </div>

            <v-form v-if="form" v-model="valid">
                <CheckFormSLOAvailability v-if="check.id === 'SLOAvailability'" :appId="appId" :form="form" />
                <CheckFormSLOLatency v-else-if="check.id === 'SLOLatency'" :appId="appId" :form="form" />
                <CheckConfigForm v-else :form="form" :check="check" :appId="appId" />

                <div v-if="check.id.startsWith('SLO')" class="my-3">
                    告警:
                    <div>
                        <ul v-if="integrations && integrations.length">
                            <li v-for="i in integrations">
                                <span>{{ i.name }}</span>
                                <span v-if="i.details" class="grey--text"> ({{ i.details }})</span>
                            </li>
                        </ul>
                        <div v-else class="grey--text">没有配置通知集成。</div>
                        <v-btn
                            color="primary"
                            small
                            :to="{ name: 'project_settings', params: { tab: 'notifications' } }"
                            :disabled="appId === '::'"
                            class="mt-1"
                        >
                            配置通知集成
                        </v-btn>
                    </div>
                </div>

                <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text class="my-3">
                    {{ error }}
                </v-alert>
                <v-alert v-if="message" color="green" outlined text class="my-3">
                    {{ message }}
                </v-alert>
                <v-btn block color="primary" @click="save" :disabled="!(valid && changed)" :loading="saving" class="mt-5">保存</v-btn>
            </v-form>
        </v-card>
    </v-dialog>
</template>

<script>
import CheckFormSLOAvailability from './CheckConfigSLOAvailabilityForm.vue';
import CheckFormSLOLatency from './CheckConfigSLOLatencyForm.vue';
import CheckConfigForm from './CheckConfigForm.vue';

export default {
    props: {
        appId: String,
        check: Object,
        value: Boolean,
    },

    components: { CheckConfigForm, CheckFormSLOAvailability, CheckFormSLOLatency },

    data() {
        return {
            loading: false,
            error: '',
            message: '',
            form: null,
            integrations: null,
            saved: '',
            saving: false,
            valid: false,
            deleting: false,
            confirmation: false,
        };
    },

    watch: {
        value() {
            if (this.value) {
                this.form = null;
                this.get();
            }
        },
    },

    computed: {
        changed() {
            return !!this.form && this.saved !== JSON.stringify(this.form);
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.$api.getInspectionConfig(this.appId, this.check.id, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    this.form = null;
                    return;
                }
                this.form = data.form;
                this.integrations = data.integrations;
                this.saved = JSON.stringify(this.form);
            });
        },
        save() {
            this.saving = true;
            this.error = '';
            this.message = '';
            this.$api.saveInspectionConfig(this.appId, this.check.id, this.form, (data, error) => {
                this.saving = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('refresh');
                this.message = 'Settings were successfully updated.';
                setTimeout(() => {
                    this.message = '';
                }, 1000);
                this.get();
            });
        },
        del() {
            this.confirmation = false;
            this.deleting = true;
            this.error = '';
            this.$api.saveInspectionConfig(this.appId, this.check.id, { configs: null }, (data, error) => {
                this.deleting = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('refresh');
                this.emitValue(false);
            });
        },
        emitValue(v) {
            this.$emit('input', v);
        },
    },
};
</script>

<style scoped></style>
