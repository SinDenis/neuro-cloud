import React from 'react';
import './App.css';
import {Route, Routes} from "react-router-dom";
import Login from "./components/Login";
import Register from "./components/Register";
import Main from "./components/Main";
import Navigation from "./components/Navigation";
import ErrorPage from "./components/ErrorPage";
import ImagePage from "./components/ImagePage";
import AudioPage from "./components/AudioPage";
import RequireAuth from "./components/AuthenticatedRoute";

const App: React.FC = () => {

  return (
    <>
      <Navigation/>
      <Routes>
        <Route index element={<Main />} />
        <Route path="main" element={<Main />} />
        <Route path="images" element={<ImagePage />} />
        <Route path="audio" element={<AudioPage /> } />
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="logout" element={<Register />} />
        <Route path="*" element={<ErrorPage />} />
      </Routes>
    </>

  )

}
export default App;
