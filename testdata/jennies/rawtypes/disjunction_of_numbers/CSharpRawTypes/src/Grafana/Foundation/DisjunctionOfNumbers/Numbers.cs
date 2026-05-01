namespace Grafana.Foundation.DisjunctionOfNumbers;


public class Numbers
{
    public long Int64;
    public double Float64;
    public float Float32;

    public Numbers()
    {
    }

    public Numbers(long int64, double float64, float float32)
    {
        this.Int64 = int64;
        this.Float64 = float64;
        this.Float32 = float32;
    }
}
