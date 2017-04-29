<template>
    <div class="uk-container uk-container-large">
        <h3>prom-metric-viewer</h3>

        <form class="uk-search uk-search-default">
            <span uk-search-icon></span>
            <input class="uk-search-input" type="search" placeholder="Search metrics" v-model="filterQuery">
        </form>

        <p>We've found {{sortedMetrics.length / 2}} metrics!</p>

        <table class="uk-table uk-table-divider uk-table-hover uk-table-small">
            <thead>
            <tr>
                <th></th>
                <th>Name</th>
                <th>Type</th>
                <th>Cardinality</th>
                <th>Help</th>
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
    form.uk-search {
        width: 300px;
    }
</style>
