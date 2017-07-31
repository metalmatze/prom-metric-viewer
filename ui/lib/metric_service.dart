import 'dart:async';
import 'dart:convert';

import 'package:angular2/angular2.dart';
import 'package:http/browser_client.dart';

import 'metric.dart';

@Injectable()
class MetricService {
  final BrowserClient _http;

  MetricService(this._http);

  Future<List<Metric>> getMetrics() async {
    final url = 'http://localhost:8888/metrics.json';

    try {
      final resp = await _http.get(url);
      return JSON.decode(resp.body)
          .map((value) => new Metric.fromJson(value))
          .toList();
    } catch (e) {
      throw _handleError(e);
    }
  }

  Future<List<RawMetric>> getMetric(String name) async {
    final url = 'http://localhost:8888/metrics.json?name=' + name;

    try {
      final resp = await _http.get(url);
      return JSON.decode(resp.body)
          .map((value) => new RawMetric.fromJSON(value))
          .toList();
    } catch (e) {
      throw _handleError(e);
    }
  }

  Exception _handleError(dynamic e) {
    return new Exception('Server error; cause: $e');
  }
}
