import 'dart:async';

import 'package:angular2/angular2.dart';

import 'metric.dart';
import 'metric_service.dart';

@Component(
  selector: 'prom-metric-viewer',
  template: '<ul><li *ngFor="let metric of metrics">{{metric.name}}</li></ul>',
  directives: const [CORE_DIRECTIVES],
  providers: const [MetricService],
)
class PromMetricViewer implements OnInit {
  final MetricService _metricService;

  List<Metric> metrics = [];
  String errorMessage;

  PromMetricViewer(this._metricService);

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
