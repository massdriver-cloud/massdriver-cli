import React from "react";
import paramsSchema from "./schema-params.json";
import connectionsSchema from "./schema-connections.json";
import artifactsSchema from "./schema-artifacts.json";
import uiSchema from "./schema-ui.json";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";

const type = "{{ .Type }}";
const Form = withTheme(MaterialUITheme);

const logParams = (rjsf, event) => {
  var json = JSON.stringify(rjsf.formData)
  console.log(json)
}

const log = (type) => console.log.bind(console, type);
const InputsTemplate = (args) => (
  <Form
    schema={paramsSchema}
    uiSchema={uiSchema}
    onChange={log("changed")}
    onSubmit={logParams}
    onError={log("errors")}
  />
);
const ConnectionsTemplate = (args) => (
  <Form
    schema={connectionsSchema}
    onChange={log("changed")}
    onSubmit={logParams}
    onError={log("errors")}
  />
);
const ArtifactsTemplate = (args) => (
  <Form
    schema={artifactsSchema}
    onChange={log("changed")}
    onSubmit={logParams}
    onError={log("errors")}
  />
);

export default {
  title: `Bundles/${type}`,
  component: InputsTemplate,
};

export const Inputs = InputsTemplate.bind({});
Inputs.args = {};

export const Connections = ConnectionsTemplate.bind({});
Connections.args = {};

export const Artifacts = ArtifactsTemplate.bind({});
Artifacts.args = {};
