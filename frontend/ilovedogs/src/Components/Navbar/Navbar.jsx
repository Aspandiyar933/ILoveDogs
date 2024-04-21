import React from 'react';
import './Navbar.css'

const Navbar = () => {
  return (
    <nav className='nav'>
        <div className='nav-logo'></div>
        <ul className='nav-menu'>
            <li><a href="/">Home</a></li>
            <li><a href="/about">About</a></li>
            <li><a href="/logIn">Log in</a></li>
            <li><a href="/sin">Sign up</a></li>
        </ul>
    </nav>
  );
}

export default Navbar;
