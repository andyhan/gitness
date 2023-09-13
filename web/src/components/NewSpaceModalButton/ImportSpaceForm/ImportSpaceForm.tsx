import React, { useState } from 'react'
import { Intent } from '@blueprintjs/core'
import * as yup from 'yup'
import { Color } from '@harnessio/design-system'
import { Button, Container, Label, Layout, FlexExpander, Formik, FormikForm, FormInput, Text } from '@harnessio/uicore'
import { Icon } from '@harnessio/icons'
import { useStrings } from 'framework/strings'
import { Organization, type ImportSpaceFormData } from 'utils/GitUtils'
import css from '../NewSpaceModalButton.module.scss'

interface ImportFormProps {
  handleSubmit: (data: ImportSpaceFormData) => void
  loading: boolean
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  hideModal: any
}

const ImportSpaceForm = (props: ImportFormProps) => {
  const { handleSubmit, loading, hideModal } = props
  const { getString } = useStrings()
  const [auth, setAuth] = useState(false)
  const [step, setStep] = useState(0)
  const [buttonLoading, setButtonLoading] = useState(false)

  const formInitialValues: ImportSpaceFormData = {
    gitProvider: '',
    username: '',
    password: '',
    name: '',
    description: '',
    organization: ''
  }
  const providers = [
    { value: 'Github', label: 'Github' },
    { value: 'Gitlab', label: 'Gitlab' }
  ]
  const validationSchemaStepOne = yup.object().shape({
    gitProvider: yup.string().trim().required(getString('importSpace.providerRequired'))
  })

  const validationSchemaStepTwo = yup.object().shape({
    organization: yup.string().trim().required(getString('importSpace.orgRequired')),
    name: yup.string().trim().required(getString('importSpace.spaceNameRequired'))
  })

  return (
    <Formik
      initialValues={formInitialValues}
      formName="importSpaceForm"
      enableReinitialize={true}
      validateOnBlur
      onSubmit={handleSubmit}>
      {formik => {
        const { values } = formik
        const handleValidationClick = async () => {
          try {
            if (step === 0) {
              await validationSchemaStepOne.validate(formik.values, { abortEarly: false })
              setStep(1)
            } else if (step === 1) {
              await validationSchemaStepTwo.validate(formik.values, { abortEarly: false })
              setButtonLoading(true)
            } // eslint-disable-next-line @typescript-eslint/no-explicit-any
          } catch (err: any) {
            formik.setErrors(
              // eslint-disable-next-line @typescript-eslint/no-explicit-any
              err.inner.reduce((acc: { [x: string]: any }, current: { path: string | number; message: string }) => {
                acc[current.path] = current.message
                return acc
              }, {})
            )
          }
        }
        const handleImport = async () => {
          console.log('ds')
          await handleSubmit(formik.values)
          console.log('pos')
          setButtonLoading(false)
        }
        return (
          <Container width={'97%'}>
            <FormikForm>
              {step === 0 ? (
                <>
                  <Container width={'70%'}>
                    <Layout.Horizontal>
                      <Icon className={css.icon} name="code-info" size={16} />
                      <Text padding={{ left: 'small' }} font={{ size: 'small' }}>
                        {getString('importSpace.content')}
                      </Text>
                    </Layout.Horizontal>
                  </Container>
                  <hr className={css.dividerContainer} />
                  <Container className={css.textContainer} width={'70%'}>
                    <FormInput.Select
                      value={{ value: values.gitProvider, label: values.gitProvider } || providers[0]}
                      name={'gitProvider'}
                      label={getString('importSpace.gitProvider')}
                      items={providers}
                      className={css.selectBox}
                    />
                    {formik.errors.gitProvider ? (
                      <Text
                        margin={{ top: 'small', bottom: 'small' }}
                        color={Color.RED_500}
                        icon="circle-cross"
                        iconProps={{ color: Color.RED_500 }}>
                        {formik.errors.gitProvider}
                      </Text>
                    ) : null}
                    <Layout.Horizontal flex>
                      {getString('importSpace.authorization')}
                      <Container padding={{ left: 'small' }} width={'100%'}>
                        <hr className={css.dividerContainer} />
                      </Container>
                    </Layout.Horizontal>
                    <FormInput.Text
                      name="username"
                      label={getString('userName')}
                      placeholder={getString('importRepo.userPlaceholder')}
                      tooltipProps={{
                        dataTooltipId: 'spaceUserTextField'
                      }}
                    />
                    {formik.errors.username ? (
                      <Text
                        margin={{ top: 'small', bottom: 'small' }}
                        color={Color.RED_500}
                        icon="circle-cross"
                        iconProps={{ color: Color.RED_500 }}>
                        {formik.errors.username}
                      </Text>
                    ) : null}
                    <FormInput.Text
                      name="password"
                      label={getString('importRepo.passToken')}
                      placeholder={getString('importRepo.passwordPlaceholder')}
                      tooltipProps={{
                        dataTooltipId: 'spacePasswordTextField'
                      }}
                      inputGroup={{ type: 'password' }}
                    />
                    {formik.errors.password ? (
                      <Text
                        margin={{ top: 'small', bottom: 'small' }}
                        color={Color.RED_500}
                        icon="circle-cross"
                        iconProps={{ color: Color.RED_500 }}>
                        {formik.errors.password}
                      </Text>
                    ) : null}
                  </Container>
                </>
              ) : null}
              {step === 1 ? (
                <>
                  <Layout.Horizontal flex>
                    {/* <Container className={css.detailsLabel}> */}
                    <Text className={css.detailsLabel} font={{ size: 'small' }} flex>
                      {getString('importSpace.details')}
                    </Text>
                    {/* </Container> */}
                    <Container padding={{ left: 'small' }} width={'100%'}>
                      <hr className={css.dividerContainer} />
                    </Container>
                  </Layout.Horizontal>
                  <Container className={css.textContainer} width={'70%'}>
                    <FormInput.Text
                      name="organization"
                      label={
                        formik.values.gitProvider === Organization.GITHUB
                          ? getString('importSpace.githubOrg')
                          : getString('importSpace.gitlabGroup')
                      }
                      placeholder={getString('importSpace.orgNamePlaceholder')}
                      tooltipProps={{
                        dataTooltipId: 'importSpaceOrgName'
                      }}
                    />
                    {formik.errors.organization ? (
                      <Text
                        margin={{ bottom: 'small' }}
                        color={Color.RED_500}
                        icon="circle-cross"
                        iconProps={{ color: Color.RED_500 }}>
                        {formik.errors.organization}
                      </Text>
                    ) : null}
                    <Layout.Horizontal>
                      <Label>{getString('importSpace.importLabel')}</Label>
                      <Icon padding={{ left: 'small' }} className={css.icon} name="code-info" size={16} />
                    </Layout.Horizontal>

                    <Container className={css.importContainer} padding={'medium'}>
                      <Layout.Horizontal>
                        <FormInput.CheckBox
                          name="repositories"
                          label={getString('pageTitle.repositories')}
                          tooltipProps={{
                            dataTooltipId: 'authorization'
                          }}
                          defaultChecked
                          onClick={() => {
                            setAuth(!auth)
                          }}
                          disabled
                          padding={{ right: 'small' }}
                          className={css.checkbox}
                        />
                        <Container padding={{ left: 'xxxlarge' }}>
                          <FormInput.CheckBox
                            name="pipelines"
                            label={getString('pageTitle.pipelines')}
                            tooltipProps={{
                              dataTooltipId: 'pipelines'
                            }}
                            onClick={() => {
                              setAuth(!auth)
                            }}
                          />
                        </Container>
                      </Layout.Horizontal>
                    </Container>
                    <Container>
                      <hr className={css.dividerContainer} />
                      <FormInput.Text
                        name="name"
                        label={getString('importSpace.spaceName')}
                        placeholder={getString('enterName')}
                        tooltipProps={{
                          dataTooltipId: 'importSpaceName'
                        }}
                      />
                      {formik.errors.name ? (
                        <Text
                          margin={{ bottom: 'small' }}
                          color={Color.RED_500}
                          icon="circle-cross"
                          iconProps={{ color: Color.RED_500 }}>
                          {formik.errors.name}
                        </Text>
                      ) : null}
                      <FormInput.Text
                        name="description"
                        label={getString('importSpace.description')}
                        placeholder={getString('importSpace.descPlaceholder')}
                        tooltipProps={{
                          dataTooltipId: 'importSpaceDesc'
                        }}
                      />
                    </Container>
                  </Container>
                </>
              ) : null}

              <hr className={css.dividerContainer} />

              <Layout.Horizontal
                spacing="small"
                padding={{ right: 'xxlarge', bottom: 'large' }}
                style={{ alignItems: 'center' }}>
                {step === 1 ? (
                  <Button
                    disabled={buttonLoading}
                    text={
                      buttonLoading ? (
                        <>
                          <Container className={css.loadingIcon} width={93.5} flex={{ alignItems: 'center' }}>
                            <Icon className={css.loadingIcon} name="steps-spinner" size={16} />
                          </Container>
                        </>
                      ) : (
                        getString('importSpace.title')
                      )
                    }
                    intent={Intent.PRIMARY}
                    onClick={() => {
                      handleValidationClick()

                      if (formik.values.name !== '' || formik.values.organization !== '') {
                        handleImport()
                        setButtonLoading(false)
                      }
                      formik.setErrors({})
                    }}
                  />
                ) : (
                  <Button
                    text={getString('importSpace.next')}
                    intent={Intent.PRIMARY}
                    onClick={() => {
                      handleValidationClick()
                      if (
                        (!formik.errors.gitProvider && formik.touched.gitProvider) ||
                        (!formik.errors.username && formik.touched.username) ||
                        (!formik.errors.password && formik.touched.password)
                      ) {
                        formik.setErrors({})
                        setStep(1)
                      }
                    }}
                  />
                )}

                <Button text={getString('cancel')} minimal onClick={hideModal} />
                <FlexExpander />

                {loading && <Icon intent={Intent.PRIMARY} name="steps-spinner" size={16} />}
              </Layout.Horizontal>
            </FormikForm>
          </Container>
        )
      }}
    </Formik>
  )
}

export default ImportSpaceForm