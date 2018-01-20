/* global publishExternalAPI, createInjector */

describe('$controller', function () {
    'use strict';

    beforeEach(function () {
        delete window.angular;
        publishExternalAPI();
    });

    it('instantiates controller functions', function () {
        var injector = createInjector(['ng']);
        var $controller = injector.get('$controller');

        function MyController() {
            this.invoked = true;
        }
        var controller = $controller(MyController);

        expect(controller).toBeDefined();
        expect(controller instanceof MyController).toBe(true);
        expect(controller.invoked).toBe(true);
    });

    it('injects dependencies to controller functions', function () {
        var injector = createInjector(['ng', function ($provide) {
            $provide.constant('aDep', 42);
        }]);
        var $controller = injector.get('$controller');

        function MyController(aDep) {
            this.theDep = aDep;
        }

        var controller = $controller(MyController);

        expect(controller.theDep).toBe(42);
    });

    it('allows injecting locals to controller functions', function () {
        var injector = createInjector(['ng']);
        var $controller = injector.get('$controller');

        function MyController(aDep) {
            this.theDep = aDep;
        }

        var controller = $controller(MyController, {aDep: 42});

        expect(controller.theDep).toBe(42);
    });

    it('allows registering controllers at config time', function () {
        function MyController() {
        }

        var injector = createInjector(['ng', function ($controllerProvider) {
            $controllerProvider.register('MyController', MyController);
        }]);
        var $controller = injector.get('$controller');

        var controller = $controller('MyController');
        expect(controller).toBeDefined();
        expect(controller instanceof MyController).toBe(true);
    });

    it('allows registering several controllers in an object', function () {
        function MyController() {
        }

        function MyOtherController() {
        }

        var injector = createInjector(['ng', function ($controllerProvider) {
            $controllerProvider.register({
                MyController: MyController,
                MyOtherController: MyOtherController
            });
        }]);
        var $controller = injector.get('$controller');

        var controller = $controller('MyController');
        var otherController = $controller('MyOtherController');
        expect(controller instanceof MyController).toBe(true);
        expect(otherController instanceof MyOtherController).toBe(true);
    });

    it('allows registering controllers through modules', function () {
        var module = angular.module('myModule', []);
        module.controller('MyController', function () {
        });

        var injector = createInjector(['ng', 'myModule']);
        var $controller = injector.get('$controller');
        var controller = $controller('MyController');

        expect(controller).toBeDefined();
    });

    it('does not normally look controllers up from window', function () {
        window.MyController = function MyController() {
        };
        var injector = createInjector(['ng']);
        var $controller = injector.get('$controller');
        expect(function () {
            $controller('MyController');
        }).toThrow();
    });

    it('looks up controllers from window if so configured', function () {
        window.MyController = function MyController() {
        };
        var injector = createInjector(['ng', function ($controllerProvider) {
            $controllerProvider.allowGlobals();
        }]);

        var $controller = injector.get('$controller');
        var controller = $controller('MyController');

        expect(controller).toBeDefined();
        expect(controller instanceof window.MyController).toBe(true);
    });

    it('can return semi-constructed controller', function () {
        var injector = createInjector(['ng']);
        var $controller = injector.get('$controller');

        function MyController() {
            this.constructed = true;
            this.myAttrWhenConstructed = this.myAttr;
        }

        var controller = $controller(MyController, null, true);
        expect(controller.constructed).toBeUndefined();
        expect(controller.instance).toBeDefined();

        controller.instance.myAttr = 42;
        var actualController = controller();

        expect(actualController.constructed).toBeDefined();
        expect(actualController.myAttrWhenConstructed).toBe(42);
    });

    it('can return semi-constructed controller when using array injection', function () {
        var injector = createInjector(['ng', function ($provide) {
            $provide.constant('aDep', 42);
        }]);
        var $controller = injector.get('$controller');

        function MyController(aDep) {
            this.aDep = aDep;
            this.constructed = true;
        }

        var controller = $controller(['aDep', MyController], null, true);
        expect(controller.constructed).toBeUndefined();
        var actualController = controller();
        expect(actualController.constructed).toBeDefined();
        expect(actualController.aDep).toBe(42);
    });

    it('can bind semi-constructed controller to scope', function () {
        var injector = createInjector(['ng']);
        var $controller = injector.get('$controller');

        function MyController(aDep) {
        }

        var scope = {};

        var controller = $controller(MyController, {$scope: scope}, true, 'myCtrl');
        expect(controller.constructed).toBeUndefined();
        expect(scope.myCtrl).toBe(controller.instance);
    });

});