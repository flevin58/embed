# embed

This command line utility was written to allow embedding assets for game development, but can be used for other purposes once you understand how it works.
Embedding all game assets allows to redistribute a standalone program.

I think the best way to explain how to use it is with an example.

Suppose you have the following project setup:

```
+assets
|-- fonts
    |-- monospace.ttf
|-- images
    |-- sprite1.png 
    |-- sprite2.png 
    |-- sprite3.png
|-- sounds 
    |-- music.ogg
    |-- laser.ogg
    |-- hit.ogg
```

The simplest way is to run the command:
```
embed assets
```

An **embed.go** file will be created inside any subdirectory that contains at least one file.

For instance the sound subfolder will contain a file as follows:

```
package sound

//go:embed music.ogg
var Music_ogg []byte

//go:embed laser.ogg
var Laser_ogg []byte

//go:embed hit.ogg
var Hit_ogg []byte
```

Note that the var name is the file name with the first letter made uppercase to make the var public. Also, the . (dot) of the extension has been replaced by an _ (underscore). This will make easy to reference the variables as follows:

```
// This example uses Raylib
laserWave = rl.LoadWaveFromMemory(sound.Laser_ogg)
laserSound = rl.LoadSoundFromWave(laserWave)
```
