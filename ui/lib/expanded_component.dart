import 'dart:async';

import 'package:angular2/angular2.dart';

import 'metric.dart';
import 'metric_service.dart';

@Component(
  selector: 'metric-expended',
  templateUrl: 'expanded_component.html',
  styleUrls: const ['expanded_component.css'],
  directives: const [CORE_DIRECTIVES],
  providers: const [MetricService],
)
class ExpandedComponent implements OnInit {
  final MetricService _metricService;

  @Input('metric')
  Metric metric;

  List<RawMetric> metrics;
  String errorMessage;

  ExpandedComponent(this._metricService);

  @override
  ngOnInit() {
    getRawMetric();
  }

  Future<Null> getRawMetric() async {
    try {
      metrics = await _metricService.getMetric(metric.name);
    } catch (e) {
      errorMessage = e.toString();
    }
  }
}
