package detect

type CostAnomalyDetector struct {
	alpha      float64
	multiplier float64
}

func NewCostAnomalyDetector(alpha, multiplier float64) *CostAnomalyDetector {
	return &CostAnomalyDetector{alpha: alpha, multiplier: multiplier}
}
func (d *CostAnomalyDetector) CheckAndUpdate(currentEWMA, newCost float64) (isAnomaly bool, updatedEWMA float64) {
	isAnomaly = currentEWMA > 0 && newCost > currentEWMA*d.multiplier
	if currentEWMA == 0 {
		updatedEWMA = newCost
	} else {
		updatedEWMA = d.alpha*newCost + (1-d.alpha)*currentEWMA
	}
	return isAnomaly, updatedEWMA
}
