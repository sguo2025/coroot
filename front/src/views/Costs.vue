<template>
    <Views :loading="loading" :error="error">
        <v-alert v-if="!loading && !error && !nodes.length" color="info" outlined text>
            根因分析 目前支持对在 AWS、GCP 和 Azure 上运行的服务进行成本监控。每个节点上的代理需要访问云元数据服务，以获取实例元数据，例如区域、可用区和实例类型。
        </v-alert>

        <v-alert v-if="custom_pricing" color="info" outlined text>
            节点所在环境不在受支持的云平台上，或者代理无法访问云元数据。<br />
            在这种情况下，将使用自定义定价，您可以调整每核vCPU和每GB内存的价格。
            <CustomCloudPricing />
        </v-alert>

        <NodesCosts v-if="nodes.length" :nodes="nodes" />
        <ApplicationsCosts v-if="applications.length" :applications="applications" />
    </Views>
</template>

<script>
import Views from '@/views/Views.vue';
import NodesCosts from '@/components/NodesCosts.vue';
import ApplicationsCosts from '@/components/ApplicationsCosts.vue';
import CustomCloudPricing from '@/components/CustomCloudPricing.vue';

export default {
    components: { Views, ApplicationsCosts, NodesCosts, CustomCloudPricing },

    data() {
        return {
            nodes: [],
            applications: [],
            loading: false,
            error: '',
            custom_pricing: false,
        };
    },

    mounted() {
        this.get();
        this.$events.watch(this, this.get, 'refresh');
    },

    methods: {
        get() {
            this.loading = true;
            this.error = '';
            this.$api.getOverview('costs', '', (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.custom_pricing = data.costs.custom_pricing;
                this.nodes = data.costs.nodes || [];
                this.applications = data.costs.applications || [];
            });
        },
    },
};
</script>
