--
-- Table structure for table `example_table`
--
CREATE TABLE IF NOT EXISTS `example_table` (
  `example_uid` bigint(20) NOT NULL,
  `example_data_1` varchar(96) NOT NULL,
  `example_data_2` varchar(36) NOT NULL,
  `example_data_3` int(11) NOT NULL,
  `example_data_4` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
--
-- Indexes for table `example_table`
--
ALTER TABLE `example_table`
  ADD PRIMARY KEY (`example_uid`);
--
-- AUTO_INCREMENT for table `example_table`
--
ALTER TABLE `example_table`
  MODIFY `example_uid` bigint(20) NOT NULL AUTO_INCREMENT;
