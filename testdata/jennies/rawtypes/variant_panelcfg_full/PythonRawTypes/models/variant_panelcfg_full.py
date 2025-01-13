import typing


class Options:
    timeseries_option: str

    def __init__(self, timeseries_option: str = ""):
        self.timeseries_option = timeseries_option

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "timeseries_option": self.timeseries_option,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "timeseries_option" in data:
            args["timeseries_option"] = data["timeseries_option"]        

        return cls(**args)


class FieldConfig:
    timeseries_field_config_option: str

    def __init__(self, timeseries_field_config_option: str = ""):
        self.timeseries_field_config_option = timeseries_field_config_option

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "timeseries_field_config_option": self.timeseries_field_config_option,
        }
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "timeseries_field_config_option" in data:
            args["timeseries_field_config_option"] = data["timeseries_field_config_option"]        

        return cls(**args)



