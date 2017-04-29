<template>
    <table class="table table-hover">
        <thead>
        <tr>
            <th>Element</th>
            <th>Value</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="metric in metrics">
            <td>{{metric.element}}</td>
            <td class="value">{{metric.value}}</td>
        </tr>
        </tbody>
    </table>
</template>

<script>
    import axios from 'axios';

    export default {
        props: [
            'name',
        ],
        data() {
            return {
                metrics: [],
            }
        },
        created() {
            this.fetchData();
        },
        methods: {
            fetchData() {
                let self = this;
                axios.get(`/metrics.json?name=${this.name}`)
                    .then(function (res) {
                        self.metrics = res.data;
                    })
                    .catch(function (err) {
                        alert(err);
                    });
            }
        }
    }
</script>

<style scoped>
    /*table {*/
        /*cursor: default;*/
        /*box-shadow: 2px 2px 5px #aaa;*/
        /*font-size: 12px;*/
        /*background-color: #f1f1f1 !important;*/
        /*border: 1px solid #e0e0e0;*/
    /*}*/

    /*table td {*/
        /*word-break: break-all;*/
    /*}*/

    /*table td.value {*/
        /*width: 256px;*/
    /*}*/
</style>
