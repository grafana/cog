import typing
from ..cog import variants as cogvariants
from ..cog import runtime as cogruntime


class Dashboard:
    title: str
    panels: typing.Optional[list['Panel']]

    def __init__(self, title: str = "", panels: typing.Optional[list['Panel']] = None):
        self.title = title
        self.panels = panels

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "title": self.title,
        }
        if self.panels is not None:
            payload["panels"] = self.panels
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "title" in data:
            args["title"] = data["title"]
        if "panels" in data:
            args["panels"] = data["panels"]        

        return cls(**args)


class DataSourceRef:
    type_val: typing.Optional[str]
    uid: typing.Optional[str]

    def __init__(self, type_val: typing.Optional[str] = None, uid: typing.Optional[str] = None):
        self.type_val = type_val
        self.uid = uid

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
        }
        if self.type_val is not None:
            payload["type"] = self.type_val
        if self.uid is not None:
            payload["uid"] = self.uid
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "type" in data:
            args["type_val"] = data["type"]
        if "uid" in data:
            args["uid"] = data["uid"]        

        return cls(**args)


class FieldConfigSource:
    defaults: typing.Optional['FieldConfig']

    def __init__(self, defaults: typing.Optional['FieldConfig'] = None):
        self.defaults = defaults

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
        }
        if self.defaults is not None:
            payload["defaults"] = self.defaults
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "defaults" in data:
            args["defaults"] = FieldConfig.from_json(data["defaults"])        

        return cls(**args)


class FieldConfig:
    unit: typing.Optional[str]
    custom: typing.Optional[object]

    def __init__(self, unit: typing.Optional[str] = None, custom: typing.Optional[object] = None):
        self.unit = unit
        self.custom = custom

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
        }
        if self.unit is not None:
            payload["unit"] = self.unit
        if self.custom is not None:
            payload["custom"] = self.custom
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "unit" in data:
            args["unit"] = data["unit"]
        if "custom" in data:
            args["custom"] = data["custom"]        

        return cls(**args)


class Panel:
    title: str
    type_val: str
    datasource: typing.Optional['DataSourceRef']
    options: typing.Optional[object]
    targets: typing.Optional[list[cogvariants.Dataquery]]
    field_config: typing.Optional['FieldConfigSource']

    def __init__(self, title: str = "", type_val: str = "", datasource: typing.Optional['DataSourceRef'] = None, options: typing.Optional[object] = None, targets: typing.Optional[list[cogvariants.Dataquery]] = None, field_config: typing.Optional['FieldConfigSource'] = None):
        self.title = title
        self.type_val = type_val
        self.datasource = datasource
        self.options = options
        self.targets = targets
        self.field_config = field_config

    def to_json(self) -> dict[str, object]:
        payload: dict[str, object] = {
            "title": self.title,
            "type": self.type_val,
        }
        if self.datasource is not None:
            payload["datasource"] = self.datasource
        if self.options is not None:
            payload["options"] = self.options
        if self.targets is not None:
            payload["targets"] = self.targets
        if self.field_config is not None:
            payload["fieldConfig"] = self.field_config
        return payload

    @classmethod
    def from_json(cls, data: dict[str, typing.Any]) -> typing.Self:
        args: dict[str, typing.Any] = {}
        
        if "title" in data:
            args["title"] = data["title"]
        if "type" in data:
            args["type_val"] = data["type"]
        if "datasource" in data:
            args["datasource"] = DataSourceRef.from_json(data["datasource"])
        if "options" in data:
            args["options"] = data["options"]
        if "targets" in data:
            args["targets"] = [cogruntime.dataquery_from_json(dataquery_json, data["datasource"]["type"] if data.get("datasource") is not None and data["datasource"].get("type", "") != "" else "") for dataquery_json in data["targets"]]
        if "fieldConfig" in data:
            args["field_config"] = FieldConfigSource.from_json(data["fieldConfig"])        

        return cls(**args)



