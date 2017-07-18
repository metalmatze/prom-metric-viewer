import 'package:angular2/angular2.dart';

@Component(
  selector: 'prom-metric-viewer',
  template: '<h1>Hello {{name}}</h1>',
)
class PromMetricViewer implements OnInit {
  String name = "Angular";

  @override
  ngOnInit() {
    this.name = "Prometheus";
    print('running PromMetricViewer');
  }
}
