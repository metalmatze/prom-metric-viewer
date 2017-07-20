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

  List<Metric> metrics = [];
  String errorMessage;
  String filterQuery;

  MetricViewer(this._metricService);

  @override
  ngOnInit() async {
    getMetrics();
  }

  Future<Null> getMetrics() async {
    try {
      metrics = await _metricService.getMetrics();
    } catch (e) {
      errorMessage = e.toString();
    }
  }
}
