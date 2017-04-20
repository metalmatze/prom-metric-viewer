<template>
    <div id="app">
        <div class="container-fluid">
            <h3>prom-metric-viewer</h3>

            <div class="form-group">
                <input type="text" class="form-control" placeholder="Text input" v-model="filterQuery">
            </div>

            <p>We've found {{sortedMetrics.length}} metrics!</p>

            <table class="table table-hover table-condensed">
                <thead>
                <tr>
                    <th @click="sortBy('name')" :class="{ active: sortKey == 'name'}">
                        Name <span class="glyphicon" :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'" v-if="sortKey === 'name'"></span>
                    </th>
                    <th @click="sortBy('type')" :class="{ active: sortKey == 'type'}">
                        Type <span class="glyphicon" :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'" v-if="sortKey === 'type'"></span>
                    </th>
                    <th @click="sortBy('cardinality')" :class="{ active: sortKey == 'cardinality'}">
                        Cardinality <span class="glyphicon" :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'" v-if="sortKey === 'cardinality'"></span>
                    </th>
                    <th @click="sortBy('help')" :class="{ active: sortKey == 'help'}">
                        Help <span class="glyphicon" :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'" v-if="sortKey === 'help'"></span>
                    </th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="metric in sortedMetrics">
                    <td>{{ metric.name }}</td>
                    <td>{{ metric.type }}</td>
                    <td>{{ metric.cardinality }}</td>
                    <td>{{ metric.help }}</td>
                </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script>
    import request from 'superagent';

    export default {
        name: 'app',
        data () {
            return {
                msg: 'Welcome to Your Vue.js App!',
                metrics: [],
                sortKey: 'name',
                sortOrder: 1,
                filterQuery: null,
            }
        },
        created() {
            this.fetchData();
        },
        computed: {
            sortedMetrics: function () {
                let data = this.metrics;

                if (this.filterQuery) {
                    data = data.slice().filter((metric) => {
                        return String(metric.name).toLowerCase().indexOf(this.filterQuery) > -1;
                    })
                }

                if (this.sortKey) {
                    data = data.slice().sort((a, b) => {
                        a = a[this.sortKey];
                        b = b[this.sortKey];
                        return (a === b ? 0 : a > b ? 1 : -1) * this.sortOrder;
                    });
                }

                return data;
            }
        },
        methods: {
            fetchData() {
                let self = this;
                request
                    .get('/metrics.json')
                    .end(function (err, res) {
                        if (err !== null || res.status !== 200) {
                            alert(err.message)
                        }
                        self.metrics = JSON.parse(res.text);
                    })
            },
            sortBy: function (key) {
                this.sortKey = key;
                this.sortOrder = this.sortOrder * -1;
            }
        }
    }
</script>

<style>
    .table>thead>tr>th.active {
        background: none;
    }
</style>
