const express = require('express');
const router = express.Router();

const authController = require('../controllers/auth');

router.get('/', authController.isLoggedIn, (req, res) => {
    res.render('index', {
        // Controlles by front-end
        user: req.user
    })
});

router.get('/register', (req, res) => {
    res.render('register')
});

router.get('/login', (req, res) => {
    res.render('login')
});

router.get('/profile', authController.isLoggedIn, (req, res) => {

    if( req.user ) {
        res.render('profile', { 
            user: req.user});
    } else {
        res.redirect('/login');
    }
    
});


module.exports = router;