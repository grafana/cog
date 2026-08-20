namespace Grafana.Foundation.Constraints;

using System.Collections.Generic;

public class SomeStruct
{
    public ulong Id;
    public ulong MaybeId;
    public ulong GreaterThanZero;
    public long Negative;
    public string Title;
    public Dictionary<string, string> Labels;
    public List<string> Tags;
    public string Regex;
    public string NegativeRegex;
    public List<string> MinMaxList;
    public List<string> UniqueList;
    public List<long> FullConstraintList;

    public SomeStruct()
    {
        this.Id = 0UL;
        this.GreaterThanZero = 0UL;
        this.Negative = 0L;
        this.Title = "";
        this.Labels = new Dictionary<string, string>();
        this.Tags = new List<string>();
        this.Regex = "";
        this.NegativeRegex = "";
        this.MinMaxList = new List<string>();
        this.UniqueList = new List<string>();
        this.FullConstraintList = new List<long>();
    }

    public SomeStruct(ulong id, ulong maybeId, ulong greaterThanZero, long negative, string title, Dictionary<string, string> labels, List<string> tags, string regex, string negativeRegex, List<string> minMaxList, List<string> uniqueList, List<long> fullConstraintList)
    {
        this.Id = id;
        this.MaybeId = maybeId;
        this.GreaterThanZero = greaterThanZero;
        this.Negative = negative;
        this.Title = title;
        this.Labels = labels;
        this.Tags = tags;
        this.Regex = regex;
        this.NegativeRegex = negativeRegex;
        this.MinMaxList = minMaxList;
        this.UniqueList = uniqueList;
        this.FullConstraintList = fullConstraintList;
    }
}
