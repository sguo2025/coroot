<template>
    <div class="mx-auto">
        <h1 class="text-h5 mb-5">配置</h1>

        <v-tabs :value="tab" height="40" show-arrows slider-size="2">
            <v-tab v-for="t in tabs" :key="t.id" :to="{ params: { tab: t.id } }" :disabled="t.disabled" :tab-value="t.id" exact>
                {{ t.name }}
            </v-tab>
        </v-tabs>

        <template v-if="!tab">
            <h2 class="text-h5 my-5">工程名称</h2>
            <ProjectSettings :projectId="projectId" />

            <template v-if="projectId">
                <h2 class="text-h5 mt-10 mb-5">状态</h2>
                <ProjectStatus :projectId="projectId" />

                <h2 class="text-h5 mt-10 mb-5">API keys</h2>
                <p></p>
                <ProjectApiKeys />

                <h2 class="text-h5 mt-10 mb-5">Danger zone</h2>
                <ProjectDelete :projectId="projectId" />
            </template>
        </template>

        <template v-if="tab === 'prometheus'">
            <h1 class="text-h5 my-5">
                Prometheus 集成
                <a href="https://docs.coroot.com/configuration/prometheus" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <IntegrationPrometheus />
        </template>

        <template v-if="tab === 'clickhouse'">
            <h1 class="text-h5 my-5">
                ClickHouse 集成
                <a href="https://docs.coroot.com/configuration/clickhouse" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <p>
                根因分析先锋存储 <a href="https://docs.coroot.com/logs" target="_blank">日志</a>, <a href="https://docs.coroot.com/tracing" target="_blank">追踪</a>,
                and <a href="https://docs.coroot.com/profiling" target="_blank">分析</a> in the ClickHouse 数据库。
            </p>
            <IntegrationClickhouse />
        </template>

        <template v-if="tab === 'ai'">
            <h1 class="text-h5 my-5">AI 驱动的根因分析</h1>
            <IntegrationAI />
        </template>

        <template v-if="tab === 'aws'">
            <h1 class="text-h5 my-5">AWS 集成</h1>
            <IntegrationAWS />
        </template>

        <template v-if="tab === 'inspections'">
            <h1 class="text-h5 my-5">
                检查配置
                <a href="https://docs.coroot.com/inspections" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <Inspections />
        </template>

        <template v-if="tab === 'applications'">
            <h2 class="text-h5 my-5" id="categories">
                应用分类
                <a href="https://docs.coroot.com/configuration/application-categories" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h2>
            <p>
                您可以通过定义 <a href="https://en.wikipedia.org/wiki/Glob_(programming)" target="_blank">glob 模式</a> 将您的应用程序组织成组。对于 Kubernetes 应用程序，分类也可以通过注解 Kubernetes 对象来定义。请参考
                <a href="https://docs.coroot.com/configuration/application-categories" target="_blank">文档</a> 了解更多详情。
            </p>
            <ApplicationCategories />

            <h2 class="text-h5 mt-10 mb-5" id="custom-applications">
                自定义应用程序
                <a href="https://docs.coroot.com/configuration/custom-applications" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h2>

            <p>根因分析先锋使用以下方法将单个容器分组为应用程序:</p>

            <ul class="mb-3">
                <li><b>Kubernetes metadata</b>: Pods 被分组为 Deployments, StatefulSets 等。</li>
                <li>
                    <b>Non-Kubernetes containers</b>: 非 Kubernetes 容器，如 Docker 容器或 Systemd 单元，按其名称分组为应用程序。例如，不同主机上名为 <var>mysql</var> 的 Systemd 服务被分组为一个名为 <var>mysql</var> 的单个应用程序。
                </li>
            </ul>

            <p>
                这种默认方法在大多数情况下效果很好。然而，由于没有人比您更了解您的系统，根因分析先锋允许您手动调整应用程序分组，以更好地满足您的特定需求。您可以通过定义 <a href="https://en.wikipedia.org/wiki/Glob_(programming)" target="_blank">glob 模式</a> 来匹配所需的应用程序实例。注意，这不适用于可以通过注解 Kubernetes 对象进行自定义的 Kubernetes 应用程序。请参考
                <a href="https://docs.coroot.com/configuration/custom-applications" target="_blank">文档</a> 了解更多详情。
            </p>

            <CustomApplications />
        </template>

        <template v-if="tab === 'notifications'">
            <h1 class="text-h5 my-5">
                通知集成
                <a href="https://docs.coroot.com/alerting/slo-monitoring" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <Integrations />
        </template>

        <template v-if="tab === 'organization'">
            <h1 class="text-h5 my-5">
                Users
                <a href="https://docs.coroot.com/configuration/authentication" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <Users />
            <h1 class="text-h5 mt-10 mb-5">
                Role-Based Access Control (RBAC)
                <a href="https://docs.coroot.com/configuration/rbac" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <RBAC />
            <h1 class="text-h5 mt-10 mb-5">
                Single Sign-On (SAML)
                <a href="https://docs.coroot.com/configuration/authentication/#single-sign-on-sso" target="_blank">
                    <v-icon>mdi-information-outline</v-icon>
                </a>
            </h1>
            <SSO />
        </template>

        <template v-if="tab === 'cloud'">
            <Cloud />
        </template>
    </div>
</template>

<script>
import ProjectSettings from './ProjectSettings.vue';
import ProjectStatus from './ProjectStatus.vue';
import ProjectApiKeys from './ProjectApiKeys.vue';
import ProjectDelete from './ProjectDelete.vue';
import Inspections from './Inspections.vue';
import ApplicationCategories from './ApplicationCategories.vue';
import Integrations from './Integrations.vue';
import IntegrationPrometheus from './IntegrationPrometheus.vue';
import IntegrationClickhouse from './IntegrationClickhouse.vue';
import IntegrationAWS from './IntegrationAWS.vue';
import CustomApplications from './CustomApplications.vue';
import Users from './Users.vue';
import RBAC from './RBAC.vue';
import SSO from './SSO.vue';
import IntegrationAI from '@/views/IntegrationAI.vue';
import Cloud from './cloud/Cloud.vue';

export default {
    props: {
        projectId: String,
        tab: String,
    },

    components: {
        IntegrationAI,
        CustomApplications,
        IntegrationPrometheus,
        IntegrationClickhouse,
        IntegrationAWS,
        Inspections,
        ProjectSettings,
        ProjectStatus,
        ProjectApiKeys,
        ProjectDelete,
        ApplicationCategories,
        Integrations,
        Users,
        RBAC,
        SSO,
        Cloud,
    },

    mounted() {
        if (!this.tabs.find((t) => t.id === this.tab)) {
            this.$router.replace({ params: { tab: undefined } });
        }
    },

    computed: {
        tabs() {
            const disabled = !this.projectId;
            let tabs = [
                { id: undefined, name: '通用' },
                { id: 'prometheus', name: 'Prometheus', disabled },
                { id: 'clickhouse', name: 'Clickhouse', disabled },
                { id: 'ai', name: 'AI' },
                { id: 'cloud', name: 'Coroot Cloud' },
                { id: 'aws', name: 'AWS', disabled },
                { id: 'inspections', name: '检查', disabled },
                { id: 'applications', name: '应用', disabled },
                { id: 'notifications', name: '通知', disabled },
                { id: 'organization', name: '组织' },
            ];
            if (this.$coroot.edition === 'Enterprise') {
                tabs = tabs.filter((t) => t.id !== 'cloud');
            }
            return tabs;
        },
    },
};
</script>

<style scoped></style>
