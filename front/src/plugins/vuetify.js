import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';


import 'vuetify/styles' // Vuetify 样式
import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'

Vue.use(Vuetify);

export default new Vuetify({
    icons: {
        defaultSet: 'mdi',
        aliases,
        sets: { mdi },
    },
    theme: {
        themes: {
            light: {
                secondary: colors.blue.lighten1,
            },
            dark: {
                secondary: colors.blue.lighten1,
            },
        },
    },
});
