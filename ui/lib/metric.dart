class Metric {
  String name;
  String help;
  String type;
  int cardinality;

  bool expanded = false;

  Metric(this.name, this.help, this.type, this.cardinality);

  factory Metric.fromJson(Map<String, dynamic> metric) {
    return new Metric((metric['name']), metric['help'], metric['type'],
        metric['cardinality']);
  }
}

class RawMetric {
  String element;
  double value;

  RawMetric(this.element, this.value);
  
  factory RawMetric.fromJSON(Map<String, dynamic> rawMetric) {
    return new RawMetric(rawMetric['element'], rawMetric['value']);
  }
}
