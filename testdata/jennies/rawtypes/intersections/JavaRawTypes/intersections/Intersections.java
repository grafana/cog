package intersections;

import externalPkg.AnotherStruct;

public class Intersections extends SomeStruct, AnotherStruct {
    public String fieldString;
    public Integer fieldInteger;
    public Intersections(String fieldString,Integer fieldInteger) {
        this.fieldString = fieldString;
        this.fieldInteger = fieldInteger;
    }
    
}
