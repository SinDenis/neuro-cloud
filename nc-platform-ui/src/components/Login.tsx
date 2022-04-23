import React, {useState} from "react";
import axios from "axios";
import Cookies from "universal-cookie";
import {useNavigate} from "react-router-dom";

interface User {
  username: string;
  password: string;
}

interface LoginResponse {
  jwt: string;
}

const Login: React.FC = () => {

  const [user, setUser] = useState<User>({
    username: '',
    password: ''
  });
  const navigate = useNavigate();
  const onChangeUsername = (username: string) => setUser({...user, username})
  const onChangePassword = (password: string) => setUser({...user, password})
  const onSubmit = () => {
    axios.post<LoginResponse>('https://neuro-cloud-back-sindenis.cloud.okteto.net/api/login', {}, { auth: {...user} })
      .then(r => {
        new Cookies().set('jwt', r.data.jwt);
        navigate('/main');
      })
      .catch(err => console.log(err))
   }

  return (
    <div className="container center">
      <div className="row">
        <form className="col s6 offset-s3">
          <div className="input-field">
            <input id="first_name" type="text" className="validate"
                   value={user?.username}
                   onChange={event => onChangeUsername(event.target.value)}
            />
            <label htmlFor="first_name">Логин</label>
          </div>
          <div className="input-field">
            <input id="last_name" type="password" className="validate"
                   value={user?.password}
                   onChange={event => onChangePassword(event.target.value)}
            />
            <label htmlFor="last_name">Пароль</label>
          </div>
        </form>
      </div>
      <button className="btn waves-effect waves-light" name="action" onClick={onSubmit}>Войти</button>
    </div>
  )
}

export default Login
