export const environment = {
    production: false,
    appPrefix: 'app_wp_dev_',
    baseHref: '/',
    administratorPrefix: 'administrator',
    locale: {
        default: 'vi',
        list: ['vi', 'en'],
        version: '20220629200022',
    },
    google: {
        clientVerify: false,
        recaptchaV3: {
            siteKey: '6LfsaN0kAAAAADrFK5oKVgn-dQ6KMezPP-Zt8b1z',
        },
    },
    rsa: {
        isActive: false,
        server: {
            publicKeyEncode:'',
        },
    },
    roles: {
        superAdmin: 1,
        admin: 2,
        dev: 3,
        user:4
    },
    user: {
        avatar: {
            default: './assets/styles/downloader/images/logo-mini.png',
        },
    },
    telegram: {
        prefix: 'https://t.me/',
        botConfig: {
            botName: ''
        }
    },
    actions: [
        {
            id: 1,
            name: 'GET'
        },
        {
            id: 2,
            name: 'POST'
        },
        {
            id: 3,
            name: 'PUT'
        },
        {
            id: 4,
            name: 'DELETE'
        }
    ],
    backendServer: {
        host: 'localhost',
        port: 4040,
        prefix: 'api/v1',
        url: 'http://localhost:4040',
        paths: {
            auth: {
                login: 'auth/login',
                logout: 'auth/logout',
                social: {
                    request: 'auth/{SOCIAL_NAME}/request',
                    callback: 'auth/{SOCIAL_NAME}/callback',
                },
                signup: 'auth/register',
                google: {
                    request: 'auth/google/request',
                    callback: 'auth/google/callback',
                },
                otpVerify: 'auth/otp/verify',
            },
            administrator: {
                role: {

                },
                user: {

                },
                route: {

                },
                permission: {

                },
                actionLog: {
                    list: 'core/action-log/list'
                },
            },
            app: {
                admin: {
                },
                user:{

                }
            }
        },
    },
    state: {
        login: {
            loginConfirm: 0,
            loginStart: 1,
            loginExcute: 2,
            loginEnd: 3,
            loginFailed: -1
        }
    }
};
