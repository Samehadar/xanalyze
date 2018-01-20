﻿ param (

# Import my custom module file. Running with -force will remove an existing module session and load fresh. 
Import-Module .\manage-package.psm1 -Force

# Command section
Write-Host ''

# Call my custom remove-package function within manage-package.psm1 
remove-packages -RootFolder $RootFolder -fileToDelete $fileToDelete -DaysOld $DaysOld -DELETE $DELETE