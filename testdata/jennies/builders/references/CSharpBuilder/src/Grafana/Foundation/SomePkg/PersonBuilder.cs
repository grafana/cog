namespace Grafana.Foundation.SomePkg;


public class PersonBuilder : Cog.IBuilder<Person>
{
    protected readonly Person @internal;

    public PersonBuilder()
    {
        this.@internal = new Person();
    }

    public PersonBuilder Name(Name name)
    {
        this.@internal.Name = name;
        return this;
    }

    public Person Build()
    {
        return this.@internal;
    }
}
