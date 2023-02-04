# <img src="https://github.com/vinegar-dev/vinegar/blob/master/desktop/vinegar.svg" width="48"> Vinegar
A transparent wrapper for Roblox Player and Roblox Studio.

# Features
+ Logging for stderr
+ Handling arguments parsing and forwarding of RobloxPlayerLauncher (to be used)
+ Custom execution of wine program within wineprefix
+ Fast finding of Roblox Player and Roblox Studio
+ Clean wine log output
+ Automatic applying of [RCO](https://github.com/L8X/Roblox-Client-Optimizer) FFlags
+ (Untested) Automatic usage of the Nvidia dedicated gpu.
+ Deletion of empty log files
+ Sets up a Wine prefix automatically
+ Automatically fetch and install Roblox Player, Studio and rbxfpsunlocker
+ Browser launch (testing)
+ Faster startup of rbxfpsunlocker and the Roblox Player

# TODO
+ Fetch latest version of Roblox, when RobloxPlayerLauncher is not used.
+ Better log names
+ Fetch latest version of rbxfpsunlocker, removes safety checksumming
+ Add watchdog for unlocker in flatpak? This needs investigation.
+ Automatically kill wineprefix when Roblox has exited
+ Automatic killing of wineprefix
+ Add installation failure detection
