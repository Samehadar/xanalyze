// ---------------------------------------------------------------------------------------------------------------------
// Brief Description of admin.
//
// @module admin
// ---------------------------------------------------------------------------------------------------------------------

function AdminControllerFactory()
{
    function AdminController() {}

    return new AdminController();
} // end AdminControllerFactory

// ---------------------------------------------------------------------------------------------------------------------

angular.module('cloud-land.controllers').controller('adminController', [
    AdminControllerFactory
]);

// ---------------------------------------------------------------------------------------------------------------------