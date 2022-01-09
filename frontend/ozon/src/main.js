import ProfileView from './views/ProfileView/ProfileView.js';
import LoginView from './views/LoginView/LoginView.js';
import SignupView from './views/SignupView/SignupView.js';
import ProductView from './views/ProductView/ProductView.js';
import ProductsView from './views/ProductsView/ProductsView';
import HeaderView from './views/HeaderView/HeaderView';
import ProfilePresenter from './presenters/ProfilePresenter';
import Router from './utils/router/Router.js';
import SignupModel from './models/SignupModel';
import SignupPresenter from './presenters/SignupPresenter';
import ProfileModel from './models/ProfileModel';
import LoginModel from './models/LoginModel';
import LoginPresenter from './presenters/LoginPresenter';
import ProductsModel from './models/ProductsModel';
import ProductsPresenter from './presenters/ProductsPresenter';
import ProductModel from './models/ProductModel';
import ProductPresenter from './presenters/ProductPresenter';
import OfflineView from './views/OfflineView/OfflineView';
import CartPresenter from './presenters/CartPresenter';
import CartView from './views/CartView/CartView';
import CartModel from './models/CartModel';
import OrderView from './views/OrderView/OrderView';
import OrderPresenter from './presenters/OrderPresenter';
import OrderModel from './models/OrderModel';
import HeaderModel from './models/HeaderModel';
import HeaderPresenter from './presenters/HeaderPresenter';
import OrdersView from './views/OrdersView/OrdersView';
import OrdersPresenter from './presenters/OrdersPresenter';
import OrdersModel from './models/OrdersModel';
import ReviewModel from './models/ReviewModel';
import ReviewPresenter from './presenters/ReviewPresenter';
import ReviewView from './views/ReviewView/ReviewView';
import WebPushModel from './models/WebPushModel';

if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/sw.js')
            .catch((registrationError) => {
                console.error('SW registration failed: ', registrationError);
            });
    });
}

const application = document.getElementById('app');

Router.root = application;

const offlineView = new OfflineView(application, null, null);

const signupPresenter = new SignupPresenter(application, SignupView, SignupModel);
const loginPresenter = new LoginPresenter(application, LoginView, LoginModel);
const profilePresenter = new ProfilePresenter(application, ProfileView, ProfileModel);
const productsPresenter = new ProductsPresenter(application, ProductsView, ProductsModel);
const productPresenter = new ProductPresenter(application, ProductView, ProductModel);
const cartPresenter = new CartPresenter(application, CartView, CartModel);
const orderPresenter = new OrderPresenter(application, OrderView, OrderModel);
const ordersPresenter = new OrdersPresenter(application, OrdersView, OrdersModel);
const reviewPresenter = new ReviewPresenter(application, ReviewView, ReviewModel);

// eslint-disable-next-line no-unused-vars
const webPushModel = new WebPushModel();

const header = document.getElementsByTagName('header')[0];
const headerPresenter = new HeaderPresenter(header, HeaderView, HeaderModel);
headerPresenter.view.show();

Router
    .register(/^\/$/, productsPresenter.view)
    .register(/^\/signup$/, signupPresenter.view)
    .register(/^\/login$/, loginPresenter.view)
    .register(/^\/profile$/, profilePresenter.view)
    .register(/^\/item(\/(?<productID>[0-9]*))?$/, productPresenter.view)
    .register(/^\/items(\/(?<category>[0-9]*)(\/(?<page>[0-9]*))?)?$/, productsPresenter.view)
    .register(/^\/search\/(?<page>[0-9]*)\/$/, productsPresenter.view)
    .register(/^\/cart$/, cartPresenter.view)
    .register(/^\/offline$/, offlineView)
    .register(/^\/order$/, orderPresenter.view)
    .register(/^\/orders(\/(?<page>[0-9]*))?$/, ordersPresenter.view)
    .register(/review$/, reviewPresenter.view)
    .register(/^\/order$/, orderPresenter.view);

Router.start();
