package default_disjunction_value;


public class DisjunctionClasses {
    protected ValueA valueA;
    protected ValueB valueB;
    protected ValueC valueC;
    protected DisjunctionClasses() {}
    public static DisjunctionClasses createValueA(ValueA valueA) {
        DisjunctionClasses disjunctionClasses = new DisjunctionClasses();
        disjunctionClasses.valueA = valueA;
        return disjunctionClasses;
    }
    public static DisjunctionClasses createValueB(ValueB valueB) {
        DisjunctionClasses disjunctionClasses = new DisjunctionClasses();
        disjunctionClasses.valueB = valueB;
        return disjunctionClasses;
    }
    public static DisjunctionClasses createValueC(ValueC valueC) {
        DisjunctionClasses disjunctionClasses = new DisjunctionClasses();
        disjunctionClasses.valueC = valueC;
        return disjunctionClasses;
    }
}
