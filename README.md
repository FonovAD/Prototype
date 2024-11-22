# Prototype - link shortener

## Problem statement
There is a need for training in optimizing Go applications. 
So I decided to create an application that will have:
1. Non-trivial business logic
2. Accessing the database/other services

That's all the requirements.\
Considering the limitations and my needs, two projects come to mind: a link shortener and my own disk.\
\+ I have a Raspberry Pi5 with a Samsung 500GB SSD 980 series lying around, on which they could be deployed.

The first implementation in line will be the link shortener.

## Project requirements

General requirements (+ must be in MVP):
1. Have a simple authorization and registration tool.
2. +When the user clicks on an abbreviated link, transfer him to the resource. 
3. +Store links for the necessary time (indefinitely) 

Performance requirements:
1. RPS > 10k
2. The ability to store up to 10e9 links
3. Service availability 0.99

Hardware:
1) Raspberry Pi5 8gb + SSD m2 Samsung 980 500GB
2) Radxa  Rock 4a + SSD m2 A DATA xpg 128GB
3) HP Mini PC 800 g1(I'll describe the details later)


## Commit design style
- feature - used when adding new application-level functionality
- fix - if you fixed some serious bug
- docs - everything related
- to the style documentation - correcting typos, correcting
- refactor formatting - refactoring the application code
- test - everything related to testing
- chore - routine code maintenance

Style:
`action (with a small letter) + for which entity + (optional details)`

Example:
`fix NoMethodError in RemoteReader`
