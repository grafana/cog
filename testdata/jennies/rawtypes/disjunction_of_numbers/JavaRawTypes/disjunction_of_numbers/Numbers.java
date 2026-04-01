package disjunction_of_numbers;


public class Numbers {
    protected Long int64;
    protected Double float64;
    protected Float float32;
    protected Numbers() {}
    public static Numbers createInt64(Long int64) {
        Numbers numbers = new Numbers();
        numbers.int64 = int64;
        return numbers;
    }
    public static Numbers createFloat64(Double float64) {
        Numbers numbers = new Numbers();
        numbers.float64 = float64;
        return numbers;
    }
    public static Numbers createFloat32(Float float32) {
        Numbers numbers = new Numbers();
        numbers.float32 = float32;
        return numbers;
    }
}
