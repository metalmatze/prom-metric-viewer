import 'dart:async';

import 'package:angular2/angular2.dart';

import 'metric.dart';
import 'metric_service.dart';

@Component(
  selector: 'metric-viewer',
  templateUrl: 'metric_viewer.html',
  directives: const [COMMON_DIRECTIVES],
  providers: const [MetricService],
)
class MetricViewer implements OnInit {
  final MetricService _metricService;

  String errorMessage;
  List<Metric> _metrics = [];
  List<Metric> sortedMetrics = [];
  String filterQuery = '';
  String sortKey = 'name';
  int sortOrder = 1;

  MetricViewer(this._metricService);

  @override
  ngOnInit() async {
    getMetrics();
  }

  Future<Null> getMetrics() async {
    try {
      _metrics = await _metricService.getMetrics();
      _setSortedMetrics();
    } catch (e) {
      errorMessage = e.toString();
    }
  }

  _setSortedMetrics() {
    if (filterQuery.isEmpty) {
      sortedMetrics = _metrics;
    } else {
      sortedMetrics = _metrics.where((metric) =>
          metric.name.toLowerCase().contains(this.filterQuery.toLowerCase())
      ).toList();
    }

    sortedMetrics.sort((a, b) {
      switch (sortKey) {
        case 'name':
          return a.name.toLowerCase().compareTo(b.name.toLowerCase()) *
              sortOrder;
        case 'type':
          return a.type.toLowerCase().compareTo(b.type.toLowerCase()) *
              sortOrder;
        case 'cardinality':
          return a.cardinality.compareTo(b.cardinality) * sortOrder;
        case 'help':
          return a.help.toLowerCase().compareTo(b.help.toLowerCase()) *
              sortOrder;
      }
    });
  }

  filter() {
    _setSortedMetrics();
  }

  sortBy(String key) {
    if (sortKey == key) {
      sortOrder = sortOrder * -1;
    } else {
      sortOrder = 1;
    }
    sortKey = key;
    _setSortedMetrics();
  }
}
