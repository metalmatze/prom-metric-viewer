<template>
    <tr v-if="metric.type=='raw' && metric.showRaw" class="raw">
        <td colspan="5">
            <table-row-raw :name="metric.name"></table-row-raw>
        </td>
    </tr>
    <tr v-else-if="metric.type!='raw'" @click="toggleRaw">
        <td>
            <span class="glyphicon"
                  :class="!metric.showRaw ? 'glyphicon-plus' : 'glyphicon-minus'">
            </span>
        </td>
        <td>{{metric.name}}</td>
        <td>{{metric.type}}</td>
        <td>{{metric.cardinality}}</td>
        <td>{{metric.help}}</td>
    </tr>
</template>

<script>
    import TableRowRaw from './TableRowRaw.vue';

    export default {
        props: [
            'metric',
        ],
        components: {
            TableRowRaw
        },
        methods: {
            toggleRaw() {
                this.$emit('toggleRaw', this.metric.name);
            }
        }
    }
</script>

<style scoped>
    .table tr td {
        overflow: hidden;
    }

    .table tr.raw {
        background: none;
    }

    .table tr.raw td {
        padding: 20px;
        cursor: default;
    }

    .table tr td > span.glyphicon {
        opacity: 0.1;
    }

    .table tr:hover {
        cursor: pointer;
    }

    .table tr:hover span.glyphicon {
        opacity: 1;
    }
</style>
