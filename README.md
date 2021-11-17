#EnekoEnergie

Hi Weave!

This is my implementation for the assignment. You might already be aware that it is not a complete implementation, yet. The past week I have been working on getting it done whenever possible, but me and the recruiter have decided that I should hand it in because you might get a glimpse of my train of thought trying to solve coding problems.

##About me
I like solving programming problems. Last year (when spare time was still a thing), I used to work on Leetcode problems. I started doing this because I noticed that I got the basics of programming down, but the real challenges for me was really in solving this kind of problems. I think I got better at it, but you be the judge. My weakness is that I think recursively about problems, which lead me thinking about edge cases too much and causes me to be stuck.

##About the implementation
EnekoEnergie is a fictitious electric/gas shipping company. They offer their services to customers in the entire country. In order to invoice customers they need data on usage that are read using metering points.

Metering points may be grouped based on:
* Metering point ID (distinct metering point)
* Metering point type ID (gas or electric metering points)

Every 15 minutes there are new readings read by these metering points. I created half of an implementation to process these readings and turn them into usage segments, which will finally be turned into total cost per metering point ID made. 

##Inputs
The inputs for the assignment are CSV (comma-seperated values) files. These are read by the built-in CSV package of Golang.

## Outputs
**I have a detailed description and visualization of the program inside the attached 'Whiteboard' file. Please refer to this for my complete thought process.**

_Disclaimer:_
The following assumptions/things should be taken into consideration whilst going through the program:
1. I have 1 week of experience in Golang. Literally, this was my first time working with it.
2. My data structures & algorithms course is due to start this year. Therefore, minimal optimizations with my limited knowledge have been made, but not more than that.
3. It's incomplete. I didn't manage to complete it yet, but I hope the basic idea can be transmitted with what's been done
4. Very little/no testing has been done. I did write a couple of unit tests for small utility functions, but that's all I managed to do.

That being said, here's a short summary of the (intended) workings of the program:
1. Read reading values in sets of 2 lines
2. Calculate the usage made in this time
3. Get the correct price for the segment based on the type and time of reading
4. If there's a invalid usage segment (negative or > 100 value):
   1. We have to either assume a linear usage based on the other calculated segments or skip it altogether.
   2. If I'd have more time, I'd probably be looking at implementing a mechanism of average usage (maybe even by time and day) and assume that value if there wasn't anything else available.
5. At the end, go through all usage segments marked as invalid and patch these by looking for any other reference usage segments we can use <- this part is neither implemented nor tested yet. I found 2 problems with it:
   1. It's inefficient to loop through all usage segments if all of them are invalid
   2. This is not working fully right now as the loop ignores patched usage segments if the first 3/4 were invalid.
6. Finally, we would sum everything up and calculate the totals


