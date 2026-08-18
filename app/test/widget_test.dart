import 'package:flutter_test/flutter_test.dart';

import 'package:petikhula_app/main.dart';

void main() {
  testWidgets('compares carton vs loose from the default inputs', (tester) async {
    await tester.pumpWidget(const PetiKhulaApp());
    await tester.tap(find.text('Compare'));
    await tester.pump();
    expect(find.textContaining('Winner:'), findsOneWidget);
  });

  test('trueCost mirrors the Go oracle: carton wins at fast off-take', () {
    final r = trueCost(
      cartonPrice: 35, loosePrice: 40, cartonSize: 24,
      spoilagePct: 2, capitalPct: 24, weekly: 24,
    );
    expect(r.winner, 'carton');
  });
}
