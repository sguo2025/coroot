<template>
    <div style="max-width: 800px">
        <p>
            根因分析先锋利用大型语言模型 (LLMs) 自动生成清晰、简洁的根因总结，帮助您的团队更快地排查问题。
        </p>
        <v-alert v-if="disabled" color="info" outlined text>
            仅在根因分析先锋企业版中提供 (每月每 CPU 核心 $1)。<br />
            <a href="https://coroot.com/account" target="_blank" class="font-weight-bold">开始</a> 您的免费试用。
        </v-alert>
        <v-alert v-if="readonly" color="primary" outlined text>
            AI 设置通过配置定义，无法通过 UI 修改。
        </v-alert>
        <v-form v-if="form" v-model="valid" :disabled="disabled || readonly" ref="form">
            <div class="subtitle-1 mt-3">模型提供商</div>
            <v-radio-group v-model="form.provider" row dense class="mt-0" hide-details>
                <v-radio value="anthropic">
                    <template #label>
                        <img :src="`${$coroot.base_path}static/img/icons/anthropic.svg`" height="20" width="20" class="mr-1" />
                        Anthropic
                    </template>
                </v-radio>
                <v-radio value="openai">
                    <template #label>
                        <img :src="`${$coroot.base_path}static/img/icons/openai.svg`" height="20" width="20" class="mr-1" />
                        OpenAI
                    </template>
                </v-radio>
                <v-radio value="openai_compatible">
                    <template #label>
                        <v-icon class="mr-1">mdi-cog-outline</v-icon>
                        OpenAI 兼容 API
                    </template>
                </v-radio>
                <v-radio value="">
                    <template #label>
                        <v-icon class="mr-1">mdi-trash-can-outline</v-icon>
                        禁用
                    </template>
                </v-radio>
            </v-radio-group>

            <template v-if="form.provider === 'anthropic'">
                <div class="subtitle-1 mt-3">API 密钥</div>
                <div class="caption">
                    要集成根因分析先锋与 Anthropic 模型，请提供您的 Anthropic API 密钥。您可以按照
                    <a href="https://docs.anthropic.com/en/api/getting-started" target="_blank">官方 Anthropic API 文档</a> 开始。
                </div>
                <v-text-field
                    v-model="form.anthropic.api_key"
                    :rules="[$validators.notEmpty]"
                    outlined
                    dense
                    hide-details
                    single-line
                    type="password"
                />
            </template>

            <template v-if="form.provider === 'openai'">
                <div class="subtitle-1 mt-3">API 密钥</div>
                <div class="caption">
                    要集成根因分析先锋与 OpenAI 模型，请提供您的 OpenAI API 密钥。请在
                    <a href="https://openai.com/index/openai-api/" target="_blank">OpenAI API 概览页面</a> 了解更多关于 API 的信息。
                </div>
                <v-text-field v-model="form.openai.api_key" :rules="[$validators.notEmpty]" outlined dense hide-details single-line type="password" />
            </template>

            <template v-if="form.provider === 'openai_compatible'">
                <div class="subtitle-1 mt-3">Base URL (基 URL) </div>
                <div class="caption">模型提供商的 API 请求的基 URL。请参考他们的文档获取配置详情。</div>
                <v-text-field v-model="form.openai_compatible.base_url" :rules="[$validators.isUrl]" outlined dense hide-details single-line />

                <div class="subtitle-1 mt-3">API 密钥</div>
                <div class="caption">要集成根因分析先锋与 OpenAI 兼容模型，请提供您的 API 密钥。</div>
                <v-text-field
                    v-model="form.openai_compatible.api_key"
                    :rules="[$validators.notEmpty]"
                    outlined
                    dense
                    hide-details
                    single-line
                    type="password"
                />

                <div class="subtitle-1 mt-3">模型</div>
                <div class="caption">要使用的模型的名称或 ID。请参考您的提供商的文档获取有效值。</div>
                <v-text-field v-model="form.openai_compatible.model" :rules="[$validators.notEmpty]" outlined dense hide-details single-line />
            </template>

            <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text class="mt-3">
                {{ error }}
            </v-alert>
            <v-alert v-if="message" color="green" outlined text class="mt-3">
                {{ message }}
            </v-alert>
            <div class="mt-3 d-flex" style="gap: 8px">
                <v-btn color="primary" @click="save" :disabled="disabled || readonly || !valid || !changed" :loading="loading">Save</v-btn>
            </div>
        </v-form>
    </div>
</template>

<script>
export default {
    components: {},

    data() {
        return {
            disabled: this.$coroot.edition !== 'Enterprise',
            readonly: false,
            form: { provider: '', anthropic: {}, openai: {}, openai_compatible: {} },
            valid: false,
            loading: false,
            error: '',
            message: '',
            saved: {},
        };
    },

    mounted() {
        this.get();
    },
    computed: {
        changed() {
            return JSON.stringify(this.form) !== JSON.stringify(this.saved);
        },
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.ai(null, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.readonly = data.readonly;
                this.form.provider = data.provider;
                this.form.anthropic = data.anthropic || {};
                this.form.openai = data.openai || {};
                this.form.openai_compatible = data.openai_compatible || {};
                this.saved = JSON.parse(JSON.stringify(this.form));
            });
        },
        save() {
            this.loading = true;
            this.error = '';
            this.message = '';
            const form = JSON.parse(JSON.stringify(this.form));
            this.$api.ai(form, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.message = 'Settings were successfully updated.';
                setTimeout(() => {
                    this.message = '';
                }, 3000);
                this.get();
            });
        },
    },
};
</script>

<style scoped></style>
