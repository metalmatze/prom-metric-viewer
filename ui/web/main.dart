// Copyright (c) 2017. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.

import 'package:angular2/angular2.dart';
import 'package:angular2/platform/browser.dart';
import 'package:angular_app/prom-metric-viewer.dart';
import 'package:http/browser_client.dart';

void main() {
  bootstrap(PromMetricViewer, [
    provide(BrowserClient, useFactory: () => new BrowserClient(), deps: [])
  ]);
}
