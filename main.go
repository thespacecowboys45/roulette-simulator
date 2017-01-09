/**
 * @date 2016.12.05
 * @author The Space Cowboy
 *
 * DESCRIPTION:
 *     A simulation for the game of American roulette.
 *
 * INDUSTRY LINKS:
 *     https://www.casinolistings.com/casinos/ecogra-approved
 *
 * INTERESTING LINKS:
 *     https://en.wikipedia.org/wiki/Fitness_proportionate_selection
 *     http://www.roulette.co.uk/guide/random-number-generation/
 *     https://www.casinolistings.com/games/rtg/free-american-roulette-game
 *     http://stackoverflow.com/questions/33003974/roulette-wheel-selection-in-genetic-algorithm-using-golang
 */
package main

import(
    "log"
    "math/rand"
    "time"
    "drandom"
)

/***********************************************
 * GLOBALS
 */
var outsideBets_Verticle1 []int = []int{1,4,7,10,13,16,19,22,25,28,31,34}
var outsideBets_Verticle2 []int = []int{2,5,8,11,14,17,20,23,26,29,32,35}
var outsideBets_Verticle3 []int = []int{3,6,9,12,15,18,21,24,27,30,33,36}
var numSpins int = 1000
var betNumber int = 1
var bets = make([][]int, 0)
var newBets = make([]BetType, 0)



var outsideBet_Verticle1, outsideBet_Verticle2, outsideBet_Verticle3 BetType

var spinStats Stats
var playerStats PlayerStats
    
type Stats struct {
    cnt_outsideBets_Verticle1 int
    cnt_outsideBets_Verticle2 int
    cnt_outsideBets_Verticle3 int
    cnt_zero int
    cnt_doubleZero int
}

// A bet type consists of:
//    1) A list (e.g. an array) of numbers considered winning number
//    2) The payout ratio (ex. 2:1 for outside Bets, 1:1 for odd/even, 35:1 straight up numbers
type BetType struct {
    numbers []int
    payoutRatio int
    amountBet int
}

type PlayerStats struct {
    betsWon int
    betsLost int
    totalAmountWon int
    totalAmountLost int
	// Variables to track consecutive wins/losses
	currentBetResult string
	previousBetResult string
    currentWinningStreak int
    currentLosingStreak int
    winningStreaks [50]int
    losingStreaks[50] int
}

/**
 * Sees if a number is in a result set
 */
func InResultSet(n int, resultSet []int) bool {
    var j int
	for j = 0; j<len(resultSet); j++ {
	    if n == resultSet[j] {
	        return true
	    }
	}    
	return false
}

func printOutBets(bets [][]int) {
    log.Printf("Printing current bets")
    
    var i int
    for i=0; i<len(bets); i++ {
        log.Printf("Bet #%d covers numbers: %+v", i+1, bets[i]) 
    }
}

func printOutNewBets(newBets []BetType) {
    log.Printf("Printing current new bets")
    
    var i int
    for i=0; i<len(newBets);i++ {
        log.Printf("New Bet #%d for $%d.00 covers numbers: %+v", i+1, newBets[i].amountBet, newBets[i].numbers)
    }
}

func seedBetType(numbers []int, payoutRatio int) BetType {
    var betType BetType
    betType.numbers = numbers
    betType.payoutRatio = payoutRatio
    betType.amountBet = 0
    return betType
}

/**
 * evaluateBets()
 *
 * @param spinResult int - the number the spin of the wheel came up with
 * @param bets [][]int - an array of an array of ints covering the numbers each bet was made upon
 *
 * @returns int, int - the # of bets won, the # of bets lost
 *
 * Looks through all the bets and examines the total number of bets won
 */
func evaluateBets(spinResult int, bets[][]int) (int, int) {
    log.Printf("Evaluating bets against spin result of: %d", spinResult)
    
    var i, betsWon, betsLost int
    
    // Look at each bet placed
    for i=0; i<len(bets); i++ {

		// See if the spinResult was in the bet placed        
        if InResultSet(spinResult, bets[i]) {
            log.Printf("Bet #%d covering numbers %+v was won!", i+1, bets[i])
            betsWon++
        } else {
            log.Printf("Bet #%d covering numbers %+v was lost.", i+1, bets[i])
            betsLost++
        }
    }
    
    return betsWon, betsLost
}

func evaluateNewBets(spinResult int, bets[]BetType) (int, int, int, int) {
    log.Printf("Evaluating new bets against spin result of: %d", spinResult)
    
	var i, betsWon, betsLost, betAmountWon, betAmountLost, totalAmountWon, totalAmountLost int
	
	// Look at each bet placed
	for i=0; i<len(bets); i++ {
	    
	    // See if the spinResult was in the bet placed
	    if InResultSet(spinResult, bets[i].numbers) {
	        betsWon++
	        betAmountWon = bets[i].amountBet * bets[i].payoutRatio
	        log.Printf("Bet #%d covering numbers %v was won for $%d.00!", i+1, bets[i].numbers, betAmountWon)
	        totalAmountWon += betAmountWon
	    } else {
	        betsLost++
	        betAmountLost = bets[i].amountBet
	        log.Printf("Bet #%d covering numbers %v was lost for $%d.00", i+1, bets[i].numbers, betAmountLost)
	        totalAmountLost += betAmountLost
	    }
	}
	return betsWon, betsLost, totalAmountWon, totalAmountLost
}

func PlaceBet(newBets []BetType, amount int, betType BetType) []BetType {
    log.Printf("PlaceBet in the amount of %d", amount)
    
    betType.amountBet = amount
    newBets = append(newBets, betType)
    
    return newBets
}

func UpdatePlayerStats(betsWon int, betsLost int, totalAmountWon int, totalAmountLost int) {
	playerStats.betsWon += betsWon
	playerStats.betsLost += betsLost
	playerStats.totalAmountWon += totalAmountWon
	playerStats.totalAmountLost += totalAmountLost
	
	var betBalance int = totalAmountWon - totalAmountLost
	
	if betBalance > 0 {
	    // Overall, the bet is considered a "win"
	    log.Printf("Overall bet considered a WIN")

	    playerStats.currentBetResult = "win"
	    
	    if playerStats.previousBetResult == "none" {
	        // This was the first bet on the wheel
	        log.Printf("\tNo previous bet made")
	    } else if playerStats.previousBetResult == "loss" {
	        // A break in the consecutive losses
	        playerStats.losingStreaks[playerStats.currentLosingStreak] += 1
	    }
	    
	    // Reset the streak 
	    playerStats.currentWinningStreak += 1
	    playerStats.currentLosingStreak = 0
	    playerStats.previousBetResult = "win"
	    log.Printf("\tCurrentWinningStreak is %d bets in a row", playerStats.currentWinningStreak)
	} else {
	    // Overall, the bet is considered a "loss"
	    log.Printf("Overall bet considered a LOSS")
	    
	    playerStats.currentBetResult = "loss"

	    if playerStats.previousBetResult == "none" {
	        // This was the first bet on the wheel
	        log.Printf("\tNo previous bet made")
	    } else if playerStats.previousBetResult == "win" {
	        playerStats.winningStreaks[playerStats.currentWinningStreak] += 1
	    }
	    
	    // Reset the streak
	    playerStats.currentWinningStreak = 0
	    playerStats.currentLosingStreak += 1
	    playerStats.previousBetResult = "loss"
	    log.Printf("\tCurrentLosingStreak is %d bets in a row", playerStats.currentLosingStreak)
	}
}

// Function is needed to track the very last bet made in the series
func FinalizePlayerConsecutiveStreaks() {
    // If the last bet made was a win
    if playerStats.currentBetResult == "win" {
        // And the previous bet made was a win
        if playerStats.previousBetResult == "win" {
            // Then increment the total count for the number of consecutive wins in a row
            playerStats.winningStreaks[playerStats.currentWinningStreak] += 1
        } else {
            // Else this is considered a first lost, so increment the number of cosecutive times player lost 1-time in a row
            playerStats.losingStreaks[1] += 1
        }
    } else {
		// If the last bet made was a loss
        // And the previous bet made was a loss
        if playerStats.previousBetResult == "loss" {
            // Then increment the total count for the number of consecutive losses in a row
            playerStats.losingStreaks[playerStats.currentLosingStreak] += 1
        } else {
            // Else this is considered a first lost, so increment the number of consecutive times player won 1-time in a row
            playerStats.winningStreaks[1] += 1
        }
    }
}

func PlaceBets() {
    log.Printf("[main][PlaceBets()][extry]")
    
    // Place bets
    bets = append(bets, outsideBets_Verticle1)
    bets = append(bets, outsideBets_Verticle3)

    printOutBets(bets)
    
    // Second stab at placing bets
    newBets = PlaceBet(newBets, 25, outsideBet_Verticle1)
    newBets = PlaceBet(newBets, 25, outsideBet_Verticle3)

    printOutNewBets(newBets)
    log.Printf("[main][PlaceBets()][exit]")
}

func main() {
    log.Printf("[main][main()][entry]")
    
    now := time.Now().Unix()
    nowNano := time.Now().UnixNano()
    rand.Seed(nowNano)
    
    // Create all potential bet types
    outsideBet_Verticle1 = seedBetType([]int{1,4,7,10,13,16,19,22,25,28,31,34}, 2)
    outsideBet_Verticle2 = seedBetType([]int{2,5,8,11,14,17,20,23,26,29,32,35}, 2)
	outsideBet_Verticle3 = seedBetType([]int{3,6,9,12,15,18,21,24,27,30,33,36}, 2)
	
    log.Printf("[main][main()]v1: %+v", outsideBet_Verticle1)
	log.Printf("[main][main()]v2: %+v", outsideBet_Verticle2)
	log.Printf("[main][main()]v3: %+v", outsideBet_Verticle3)
	
/*    
    // Place bets
    var bets = make([][]int, 0)
    bets = append(bets, outsideBets_Verticle1)
    bets = append(bets, outsideBets_Verticle3)

    printOutBets(bets)

    // Second stab at placing bets
    var newBets = make([]BetType, 0)
    newBets = PlaceBet(newBets, 25, outsideBet_Verticle1)
    newBets = PlaceBet(newBets, 25, outsideBet_Verticle3)

    printOutNewBets(newBets)
*/    
    
    // Initialize player stats
    playerStats.currentBetResult = "none"
    playerStats.previousBetResult = "none"

    log.Printf("[main][main()]v1 now: %+v", outsideBet_Verticle1)
    log.Printf("[main][main()]v3 now: %+v", outsideBet_Verticle3)
        
    log.Printf("Now:%d NowNano:%d", now, nowNano)
    for betNumber <= numSpins {
        log.Printf("==================================================")
        log.Printf("[main][main()][%d] Spinning the wheel", betNumber)
    
    	// Uhm...Er...Place your bets!    
		PlaceBets()
        
        // Spin the wheel
        //result := rand.Intn(38)
        result := drandom.Intn(38)
       	log.Printf("[main][main()][%d] result=%d", betNumber, result) 

       if result == 0 {
           log.Printf("0 found: %d", result)
           spinStats.cnt_zero++
       }

       if result == 37 {
           log.Printf("00 found: %d", result)
           spinStats.cnt_doubleZero++
       }
       
		// check if result is in outsideBets_Verticle1
		if InResultSet(result, outsideBets_Verticle1) {
		    log.Printf("result %d is in outsideBets_Verticle1", result)
		    spinStats.cnt_outsideBets_Verticle1++
		}

		if InResultSet(result, outsideBets_Verticle2) {
		    log.Printf("result %d is in outsideBets_Verticle2", result)
		    spinStats.cnt_outsideBets_Verticle2++
		}

		if InResultSet(result, outsideBets_Verticle3) {
		    log.Printf("result %d is in outsideBets_Verticle3", result)
		    spinStats.cnt_outsideBets_Verticle3++
		}
		
		betsWon, betsLost := evaluateBets(result, bets)
		log.Printf("Player won %d and lost %d bets.", betsWon, betsLost)
		
		betsWon, betsLost, totalAmountWon, totalAmountLost := evaluateNewBets(result, newBets)
		log.Printf("Player won %d and lost %d new bets.", betsWon, betsLost)
		log.Printf("Player won $%d.00 and lost $%d.00.", totalAmountWon, totalAmountLost)
		
		UpdatePlayerStats(betsWon, betsLost, totalAmountWon, totalAmountLost)
		//playerStats.betsWon += betsWon
		//playerStats.betsLost += betsLost
		//playerStats.totalAmountWon += totalAmountWon
		//playerStats.totalAmountLost += totalAmountLost
		
		
		
		// next Spin
		betNumber += 1
    }
    
    FinalizePlayerConsecutiveStreaks()
    
    log.Printf("Total Spins: %d", numSpins)
    log.Printf("Spin Stats: %+v", spinStats)
    log.Printf("Player Stats: %+v", playerStats)
}