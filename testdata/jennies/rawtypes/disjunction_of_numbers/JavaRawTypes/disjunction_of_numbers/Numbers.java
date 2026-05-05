package disjunction_of_numbers;

import java.util.Objects;

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

    @Override
    public boolean equals(Object other) {
        if (this == other) return true;
        if (!(other instanceof Numbers)) return false;
        Numbers o = (Numbers) other;
        if (!Objects.equals(this.int64, o.int64)) return false;
        if (!Objects.equals(this.float64, o.float64)) return false;
        if (!Objects.equals(this.float32, o.float32)) return false;
        return true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.int64, this.float64, this.float32);
    }
}
