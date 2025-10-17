<template>
    <div style="max-width: 800px">
        <div class="d-block d-md-flex align-center">
            <div class="flex-grow-1">
                <div><b>删除工程</b></div>
                <div>删除工程后，无法恢复。请谨慎操作。</div>
            </div>
            <div>
                <v-btn block @click="dialog = true" :disabled="readonly" color="red" outlined>删除工程</v-btn>
            </div>
        </div>
        <v-dialog v-model="dialog" max-width="600">
            <v-card v-if="loading" class="pa-10">
                <v-progress-linear indeterminate />
            </v-card>
            <v-card v-else class="pa-4">
                <div class="d-flex align-center font-weight-bold mb-4">
                    确定要删除吗？
                    <v-spacer />
                    <v-btn icon @click="dialog = false"><v-icon>mdi-close</v-icon></v-btn>
                </div>
                <p>
                    此操作无法撤销。将永久删除 <b>{{ name }}</b> 工程。
                </p>
                <p>
                    请输入 <b>{{ name }}</b> 确认
                </p>
                <v-text-field v-model="confirmation" outlined dense></v-text-field>
                <v-alert v-if="error" color="red" icon="mdi-alert-octagon-outline" outlined text>
                    {{ error }}
                </v-alert>
                <v-btn block color="red" outlined :disabled="confirmation !== name" @click="del">
                    <template v-if="$vuetify.breakpoint.mdAndUp"> 我理解后果，删除工程 </template>
                    <template v-else> 删除工程 </template>
                </v-btn>
            </v-card>
        </v-dialog>
    </div>
</template>

<script>
export default {
    props: {
        projectId: String,
    },

    data() {
        return {
            readonly: false,
            name: '',
            dialog: false,
            loading: false,
            confirmation: '',
            error: '',
        };
    },

    mounted() {
        this.get();
    },

    watch: {
        dialog(v) {
            this.confirmation = '';
            v && this.get();
        },
    },

    methods: {
        get() {
            this.error = '';
            this.loading = true;
            this.$api.getProject(this.projectId, (data, error) => {
                this.loading = false;
                if (error) {
                    this.error = error;
                    return;
                }
                this.readonly = data.readonly;
                this.name = data.name;
            });
        },
        del() {
            this.error = '';
            this.$api.delProject(this.projectId, (data, error) => {
                if (error) {
                    this.error = error;
                    return;
                }
                this.$events.emit('projects');
                this.$router.push({ name: 'index' });
            });
        },
    },
};
</script>

<style scoped></style>
