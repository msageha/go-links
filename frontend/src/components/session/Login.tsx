import { createStyles, makeStyles } from '@material-ui/core'
import React, { useState } from 'react'
import { useDispatch } from 'react-redux'
import Router from 'next/router'
import { GoogleLogin, GoogleLoginResponse, GoogleLoginResponseOffline } from "react-google-login"
import { Alert } from '@material-ui/lab'
import {getGoogleLoginApiClientId} from "../../config/constants";

const useStyles = makeStyles(theme =>
    createStyles({
        button: {
            margin: 'auto',
        },
    })
)

const Login: React.FC = () => {
    const classes = useStyles()
    const dispatch = useDispatch()

    const [loginFailed, setLoginFailed] = useState<boolean>();

    const onSuccess = (response: GoogleLoginResponse | GoogleLoginResponseOffline) => {
        const tokenObj = (response as GoogleLoginResponse).tokenObj
        if (!tokenObj) return
        dispatch(

        )
    }

    const onFailure = (error: unknown) => {
        console.error(error)
        setLoginFailed(true)
    }

    return (
        <>
        {loginFailed &&
        <Alert  severity="error">Could not sign in! Please try again.</Alert>
        }
    <GoogleLogin
            className={classes.button}
            clientId={getGoogleLoginApiClientId()}
            buttonText={"Login"}
            onSuccess={onSuccess}
            onFailure={onFailure}
            cookiePolicy={"single_host_origin"}
            />
        </>
    )
}

export default Login
