import 'package:flutter/material.dart';

void main() => runApp(const PetiKhulaApp());

/// PetiKhula — carton (peti) vs loose (khula). Nets spoilage and capital-holding
/// cost into a true cost per usable unit, mirroring the Go middleware.
class PetiKhulaApp extends StatelessWidget {
  const PetiKhulaApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'PetiKhula',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF9A5B2E), useMaterial3: true),
        home: const HomePage(),
      );
}

class Result {
  final double cartonCost, looseCost, holding, weeksToSell, breakEven;
  final String winner;
  const Result(this.cartonCost, this.looseCost, this.holding, this.weeksToSell,
      this.winner, this.breakEven);
}

/// trueCost mirrors backend/cost.go exactly.
Result trueCost({
  required double cartonPrice,
  required double loosePrice,
  required int cartonSize,
  required double spoilagePct,
  required double capitalPct,
  required double weekly,
}) {
  final spoilageFactor = 1 - spoilagePct / 100;
  final cartonSpoilageAdj = cartonPrice / spoilageFactor;
  final double weeksToSell = weekly <= 0 ? 0.0 : cartonSize / weekly;
  final holding = cartonPrice * (capitalPct / 100) * (weeksToSell / 52) / 2;
  final cartonCost = cartonSpoilageAdj + holding;
  final winner = cartonCost < loosePrice
      ? 'carton'
      : (loosePrice < cartonCost ? 'loose' : 'tie');
  final k = cartonPrice * (capitalPct / 100) * (cartonSize / 52) / 2;
  final gap = loosePrice - cartonSpoilageAdj;
  final breakEven = (gap > 0 && k > 0) ? k / gap : 0.0;
  return Result(cartonCost, loosePrice, holding, weeksToSell, winner, breakEven);
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _cp = TextEditingController(text: '35');
  final _lp = TextEditingController(text: '40');
  final _cs = TextEditingController(text: '24');
  final _sp = TextEditingController(text: '5');
  final _cap = TextEditingController(text: '24');
  final _wk = TextEditingController(text: '6');
  Result? _r;

  double _n(TextEditingController c) => double.tryParse(c.text.trim()) ?? 0;
  int _i(TextEditingController c) => int.tryParse(c.text.trim()) ?? 0;

  void _calc() => setState(() => _r = trueCost(
        cartonPrice: _n(_cp), loosePrice: _n(_lp), cartonSize: _i(_cs),
        spoilagePct: _n(_sp), capitalPct: _n(_cap), weekly: _n(_wk),
      ));

  @override
  Widget build(BuildContext context) {
    String money(double v) => '₹${v.toStringAsFixed(2)}';
    return Scaffold(
      appBar: AppBar(
        title: const Text('PetiKhula · carton vs loose'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: ListView(padding: const EdgeInsets.all(16), children: [
        const Text('Should you buy the sealed carton or loose pieces?',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
        const SizedBox(height: 12),
        Row(children: [
          Expanded(child: _f(_cp, 'Carton price/unit ₹')),
          const SizedBox(width: 12),
          Expanded(child: _f(_lp, 'Loose price/unit ₹')),
        ]),
        Row(children: [
          Expanded(child: _f(_cs, 'Units per carton')),
          const SizedBox(width: 12),
          Expanded(child: _f(_wk, 'Sell-through /week')),
        ]),
        Row(children: [
          Expanded(child: _f(_sp, 'Spoilage %')),
          const SizedBox(width: 12),
          Expanded(child: _f(_cap, 'Capital cost %/yr')),
        ]),
        const SizedBox(height: 8),
        FilledButton.icon(onPressed: _calc, icon: const Icon(Icons.compare_arrows), label: const Text('Compare')),
        const SizedBox(height: 20),
        if (_r != null) _card(money),
      ]),
    );
  }

  Widget _f(TextEditingController c, String label) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: TextField(
          controller: c,
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
          decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
          onChanged: (_) => _calc(),
        ),
      );

  Widget _card(String Function(double) money) {
    final r = _r!;
    return Card(
      color: Theme.of(context).colorScheme.primaryContainer,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text('Winner: ${r.winner.toUpperCase()}',
              style: const TextStyle(fontSize: 26, fontWeight: FontWeight.bold)),
          const Divider(height: 20),
          _row('Carton true cost/usable unit', money(r.cartonCost)),
          _row('Loose cost/usable unit', money(r.looseCost)),
          _row('Holding cost/unit', money(r.holding)),
          _row('Weeks to sell a carton', r.weeksToSell.toStringAsFixed(1)),
          _row('Break-even off-take/week',
              r.breakEven > 0 ? '${r.breakEven.toStringAsFixed(1)} units' : '—'),
        ]),
      ),
    );
  }

  Widget _row(String k, String v) => Padding(
        padding: const EdgeInsets.symmetric(vertical: 2),
        child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [Flexible(child: Text(k)), Text(v, style: const TextStyle(fontWeight: FontWeight.w600))]),
      );
}
