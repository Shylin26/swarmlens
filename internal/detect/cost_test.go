package detect

import "testing"

func TestCostAnomalyDetector_FirstEventNeverFlagged(t *testing.T) {
	d := NewCostAnomalyDetector(0.3, 5.0)
	isAnomaly, emwa := d.CheckAndUpdate(0, 0.02)
	if isAnomaly {
		t.Error("expected first event to never be flagged as anomaly, got true")
	}
	if emwa != 0.02 {
		t.Errorf("expected EWMA to seed at 0.02, got %v", emwa)
	}
}
func TestCostAnomalyDetector_NormalCostNotFlagged(t *testing.T) {
	d := NewCostAnomalyDetector(0.3, 5.0)
	ewma := 0.02
	costs := []float64{0.021, 0.019, 0.022, 0.018}
	for _, cost := range costs {
		isAnomaly, updated := d.CheckAndUpdate(ewma, cost)
		if isAnomaly {
			t.Errorf("expected cost %v close to baseline %v to not be flagged", cost, ewma)
		}
		ewma = updated
	}
}
func TestCostAnomalyDetector_SpikeIsFlagged(t *testing.T) {
	d := NewCostAnomalyDetector(0.3, 5.0)
	ewma := 0.02
	isAnomaly, _ := d.CheckAndUpdate(ewma, 0.02*6)
	if !isAnomaly {
		t.Error("expected a 6x cost spike to be flagged as anomaly, got false")
	}

}
func TestCostAnomalyDetector_EWMASmoothsTowardNewValues(t *testing.T) {
	d := NewCostAnomalyDetector(0.5, 5.0)

	_, ewma := d.CheckAndUpdate(0.10, 0.20)

	expected := 0.15
	tolerance := 0.0000001
	if diff := ewma - expected; diff < -tolerance || diff > tolerance {
		t.Errorf("expected EWMA close to %v, got %v", expected, ewma)
	}
}
