# Hand Rankings & Tie break logic documentation

In brackets is the sorted output of each 'detection function'.
In order to allow for simple tie break logic of comparing cards of both hands at index 0 -> 5.
Lower numbers are of higher rank. e.g. RankOf(R1) > RankOf(R2)

## Royal Flush
Draw

[R1 R2 R3 R4 R5]

## Straight Flush
Highest Card
Draw

[S1 S2 S3 S4 S5]

## Four of a Kind
Highest Four of a Kind
Kicker
Draw

[F1 F2 F3 F4 K]

## Full House
Highest Three of a Kind
Highest Pair
Draw

[S1 S2 S3 T1 T2]

## Flush
Highest Ranking Card
Next Highest
Next Highest
Next Highest
Next Highest
Draw

[F1 F2 F3 F4 F5]

## Straight
Highest Card
Draw

[S1 S2 S3 S4 S5]

## Three of a Kind
Highest Three of a Kind
Draw

[T1 T2 T3 K1 K2]

## Two Pair
Highest Pair
Second Highest Pair
Kicker
Draw

[P1 P2 T1 T2 K1]

## One Pair
Highest Pair
Kicker
Next Highest
Next Highest
Next Highest
Draw

[P1 P2 K1 K2 K3]

## High Card
High Card
Next Highest
Next Highest
Next Highest
Next Highest
Draw

[H1 K1 K2 K3 K4 K5]

 ## Notes
For flush and high card every card acts as a kicker