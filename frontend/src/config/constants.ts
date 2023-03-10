import getConfig from 'next/config'

export const URL_AFTER_ADMIN_LOGIN = '/dashboard'

const { serverRuntimeConfig, publicRuntimeConfig } = getConfig()

export const getGoogleLoginApiClientId = () => {
    switch (publicRuntimeConfig.EXTERNAL_API_ENV) {
        case 'production':
            return ''
        default:
            return '641170215037-aj930er4qdo10340qmjbe29lqu1274dp.apps.googleusercontent.com'
    }
}
