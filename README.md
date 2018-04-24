# simple-rpg
An rpg made in golang using the pixel library

improvement ideas:
* sound volume should be adjustable in the config file
* support of multiple levels instead of a single level
* should the player be included in the entity package since it shares so many properties already? it could be like the spawner,
    a superset of entity 
* a cool ui for the player
* rename entities "actors?"
* use the command pattern for controller events
* use more events in general, register for events and what not.
* perhaps I should make "engines" owned by a central game object (higher that world) and these
engines are responsible for managing messages between components. for example, gui should hold a reference
to the window class, but should instead register themselves with the relavent events they would like to know about.