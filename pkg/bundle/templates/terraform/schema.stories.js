import paramsSchema from './schema-params.json'
import connectionsSchema from './schema-connections.json'
import artifactsSchema from './schema-artifacts.json'
import uiSchema from './schema-ui.json'
import guideMdxString from './operator.mdx'
import Form from 'components/Form'
import GuideRenderer from 'components/GuideRenderer'

const type = '{{ .Name }}'

/**
 * This is the scope object that is given to the MDX guide when parsing and
 * rendering. Any data you are expecting from your bundle on the front end can
 * be mocked here so you can properly test your MDX guide and it's content.
 */
const guideData = {
  type
}

/**
 * These are the stories for the forms and the guide. They are generated so
 * there is no need to edit anything below this point.
 */
const logParams = (rjsf, event) => {
  var json = JSON.stringify(rjsf.formData)
  console.log(json)
}

const log = type => console.log.bind(console, type)

const FormTemplate = args => <Form {...args} />

export default {
  title: `Bundles/${type}`,
  component: FormTemplate
}

export const Inputs = FormTemplate.bind({})
Inputs.args = {
  schema: paramsSchema,
  uiSchema: uiSchema,
  onChange: log('changed'),
  onSubmit: logParams,
  onError: log('errors')
}

export const Connections = FormTemplate.bind({})
Connections.args = {
  schema: connectionsSchema,
  onChange: log('changed'),
  onSubmit: logParams,
  onError: log('errors')
}

export const Artifacts = FormTemplate.bind({})
Artifacts.args = {
  schema: artifactsSchema,
  onChange: log('changed'),
  onSubmit: logParams,
  onError: log('errors')
}

export const DynamicGuide = () => (
  <GuideRenderer mdxString={guideMdxString} scope={guideData} />
)
