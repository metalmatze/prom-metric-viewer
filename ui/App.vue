<template>
    <div id="app">
        <div class="container-fluid">
            <h3>prom-metric-viewer</h3>

            <div class="form-group">
                <input type="text" class="form-control" placeholder="Search metric names" v-model="filterQuery">
            </div>

            <p>We've found {{sortedMetrics.length / 2}} metrics!</p>

            <table class="table table-hover table-condensed">
                <thead>
                <tr>
                    <th></th>
                    <th @click="sortBy('name')" :class="{ active: sortKey == 'name'}">
                        Name <span class="glyphicon"
                                   :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'"
                                   v-if="sortKey === 'name'"></span>
                    </th>
                    <th @click="sortBy('type')" :class="{ active: sortKey == 'type'}">
                        Type <span class="glyphicon"
                                   :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'"
                                   v-if="sortKey === 'type'"></span>
                    </th>
                    <th @click="sortBy('cardinality')" :class="{ active: sortKey == 'cardinality'}">
                        Cardinality <span class="glyphicon"
                                          :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'"
                                          v-if="sortKey === 'cardinality'"></span>
                    </th>
                    <th @click="sortBy('help')" :class="{ active: sortKey == 'help'}">
                        Help <span class="glyphicon"
                                   :class="sortOrder > 0 ? 'glyphicon-menu-up' : 'glyphicon-menu-down'"
                                   v-if="sortKey === 'help'"></span>
                    </th>
                </tr>
                </thead>
                <tbody>
                <table-row
                        v-for="metric in sortedMetrics"
                        :key="metric.key"
                        :metric="metric"
                        @toggleRaw="toggleRaw">
                </table-row>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script>
    import axios from 'axios';
    import TableRow from './TableRow.vue';

    export default {
        components: {
            'TableRow': TableRow,
        },
        data () {
            return {
                metrics: [],
                sortKey: 'name',
                sortOrder: 1,
                filterQuery: null,
                raws: [],
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

                let data2 = [];
                let i;
                for (i = 0; i < data.length * 2; i++) {
                    if (i % 2 === 0) {
                        let element = Object.assign({}, data[Math.floor(i / 2)]);
                        element.key = i;
                        element.showRaw = this.raws.includes(element.name);
                        data2[i] = element;
                    } else {
                        let element = Object.assign({}, data[Math.floor(i / 2)]);
                        element.type = 'raw';
                        element.key = i;
                        element.showRaw = this.raws.includes(element.name);
                        data2[i] = element;
                    }
                }
                return data2;
            }
        },
        methods: {
            fetchData() {
                let self = this;
                axios.get('/metrics.json')
                    .then(function (res) {
                        self.metrics = res.data;
                    })
                    .catch(function (err) {
                        alert(err);
                    });
            },
            sortBy: function (key) {
                this.sortKey = key;
                this.sortOrder = this.sortOrder * -1;
            },
            toggleRaw: function (name) {
                if (this.raws.includes(name)) {
                    let index = this.raws.indexOf(name);
                    if (index > -1) {
                        this.raws.splice(index, 1);
                    }
                } else {
                    this.raws.push(name);
                }
            },
        }
    }
</script>

<style scoped>
    .table {
        font-size: 12px;
    }

    .table > thead > tr > th.active {
        background: none;
    }
</style>
