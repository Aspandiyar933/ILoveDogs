import { useEffect, useState } from 'react'
import './LogIn.css' 


const LogIn = () => { 
    const [data, setData] = useState(null)
    const [loading, setLoading] = useState(null)

    useEffect(() => {
        fetch('')
    })
  return (
  <>
    <div className="under-container b-radius-30">
        <div className="wrapper-c b-radius-30">
            <div className="container">
                <header className="pos-sticky pd-30 display-flex-sb ">
                    <a href="./MainPage.html" className=" logo"><img src="./Assets/logo.svg" alt=""/></a>
                    <div className="navigation display-flex-standart">
                        <a href="#">Catalog </a>
                        <a href="#">About us</a>
                        <a href="#">Contacts</a>
                    </div>
                </header>
                <div className="content display-flex">
                    <div className="login-bg"></div>
                    <div className="login-form-wrapper">
                         <form action="">
                            <div className="element mb-20">
                                <h4>Username</h4>
                                <input type="text" placeholder=""/>
                            </div>
                            <div className="element mb-30">
                                <h4>Password</h4>
                                <input type="password" placeholder=""/>
                            </div>
                            <div className="element display-flex mb-30">
                                <p>Don’t have a account,</p><a href="./SignUp.html">Sign up</a>
                            </div>
                            <div className="element">
                                <button>Sign In</button>
                            </div>
                         </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
    </>
    )
}

export default LogIn
